package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const mockSHA = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

func TestIsSupported(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"Dockerfile", true},
		{"dockerfile", true},
		{"my.Dockerfile", true},
		{"Dockerfile.dev", true},
		{"docker-compose.yml", true},
		{"docker-compose.yaml", true},
		{"compose.yml", true},
		{"compose.yaml", true},
		{"README.md", false},
		{"config.json", false},
		{"main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isSupported(tt.path); got != tt.want {
				t.Errorf("isSupported(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsSHA(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"nginx@sha256:" + mockSHA, true},
		{"nginx:latest", false},
		{"nginx@sha256:short", false},
		{"nginx", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := isSHA(tt.input); got != tt.want {
				t.Errorf("isSHA(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseImage(t *testing.T) {
	tests := []struct {
		ref      string
		wantRepo string
		wantTag  string
	}{
		{"nginx:latest", "nginx", "latest"},
		{"nginx", "nginx", "latest"},
		{"my.registry.com/repo:1.2.3", "my.registry.com/repo", "1.2.3"},
		{"ubuntu:20.04@sha256:xxxx", "ubuntu", "20.04"},
	}

	for _, tt := range tests {
		t.Run(tt.ref, func(t *testing.T) {
			gotRepo, gotTag := parseImage(tt.ref)
			if gotRepo != tt.wantRepo {
				t.Errorf("parseImage repo = %v, want %v", gotRepo, tt.wantRepo)
			}
			if gotTag != tt.wantTag {
				t.Errorf("parseImage tag = %v, want %v", gotTag, tt.wantTag)
			}
		})
	}
}

func TestReplaceInLine(t *testing.T) {
	digestCache["nginx:latest"] = "sha256:" + mockSHA
	digestCache["postgres:14"] = "sha256:" + mockSHA

	defer func() {
		digestCache = map[string]string{}
	}()

	tests := []struct {
		name    string
		line    string
		want    string
		wantErr bool
	}{
		{
			name: "Dockerfile FROM simple",
			line: "FROM nginx:latest",
			want: "FROM nginx:latest@sha256:" + mockSHA,
		},
		{
			name: "Dockerfile FROM with AS",
			line: "FROM nginx:latest AS builder",
			want: "FROM nginx:latest@sha256:" + mockSHA + " AS builder",
		},
		{
			name: "Docker Compose image",
			line: "    image: postgres:14",
			want: "    image: postgres:14@sha256:" + mockSHA,
		},
		{
			name: "With Quotes",
			line: `    image: "postgres:14"`,
			want: `    image: "postgres:14@sha256:` + mockSHA + `"`,
		},
		{
			name: "Already Pinned (Should ignore)",
			line: "FROM nginx:latest@sha256:" + mockSHA,
			want: "FROM nginx:latest@sha256:" + mockSHA,
		},
		{
			name: "Already Pinned with comment (Should strip comment but keep pinned)",
			line: "FROM nginx:latest@sha256:" + mockSHA + " # latest",
			want: "FROM nginx:latest@sha256:" + mockSHA + " # latest",
		},
		{
			name: "Commented out (Should ignore)",
			line: "# FROM nginx:latest",
			want: "# FROM nginx:latest",
		},
		{
			name: "Variables (Should ignore)",
			line: "FROM nginx:${VERSION}",
			want: "FROM nginx:${VERSION}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replaceInLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceInLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("replaceInLine()\n got:  %v\n want: %v", got, tt.want)
			}
		})
	}
}

func TestProcessIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pinker-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) //nolint:errcheck

	digestCache["alpine:3.18"] = "sha256:" + mockSHA

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	content := `
# This is a comment
FROM alpine:3.18
RUN echo "hello"
FROM alpine:3.18 AS builder
`
	if err := os.WriteFile(dockerfilePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	changed, err := process(dockerfilePath, dockerfilePath)
	if err != nil {
		t.Fatalf("process() failed: %v", err)
	}

	if !changed {
		t.Error("process() should have returned changed=true")
	}

	newContentBytes, err := os.ReadFile(dockerfilePath)
	if err != nil {
		t.Fatal(err)
	}
	newContent := string(newContentBytes)

	expectedSnippet := "FROM alpine:3.18@sha256:" + mockSHA

	if !strings.Contains(newContent, expectedSnippet) {
		t.Errorf("File content did not update correctly.\nGot:\n%s\nExpected to contain:\n%s", newContent, expectedSnippet)
	}
}
