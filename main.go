package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/caarlos0/log"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

const (
	shaLen = 64
)

var (
	digestCache = map[string]string{}
	imageRegex  = regexp.MustCompile(`(?i)^(FROM\s+|[\s-]*image:\s+)(["']?)([^"'\s]+)(["']?)(.*)$`)
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	log.WithField("dir", dir).Info("pinning")

	var changed, total int
	if err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !isSupported(path) {
			return nil
		}

		total++
		didChange, err := process(path, path)
		if err != nil {
			log.WithError(err).
				WithField("file", path).
				Error("could not process")
			return err
		}

		if didChange {
			changed++
			log.WithField("file", path).
				Info("updated")
		}
		return nil
	}); err != nil {
		os.Exit(1)
	}

	log.WithField("dir", dir).
		WithField("total", total).
		WithField("changed", changed).
		Info("done!")
}

func process(inPath, outPath string) (bool, error) {
	f, err := os.Open(inPath)
	if err != nil {
		return false, fmt.Errorf("process: %w", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.WithError(err).WithField("file", inPath).Warn("failed to close file")
		}
	}()

	var out strings.Builder
	s := bufio.NewScanner(f)
	changed := false

	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(strings.TrimSpace(line), "#") ||
			(!strings.Contains(strings.ToUpper(line), "FROM") && !strings.Contains(line, "image:")) {
			out.WriteString(line)
			out.WriteByte('\n')
			continue
		}

		newLine, err := replaceInLine(line)
		if err != nil {
			log.WithError(err).WithField("line", strings.TrimSpace(line)).Debug("skipping line")
			out.WriteString(line)
			out.WriteByte('\n')
			continue
		}

		if newLine != line {
			changed = true
		}
		out.WriteString(newLine)
		out.WriteByte('\n')
	}

	if err := s.Err(); err != nil {
		return false, fmt.Errorf("process: %w", err)
	}

	if changed {
		if err := os.WriteFile(outPath, []byte(out.String()), 0o644); err != nil {
			return false, fmt.Errorf("process: %w", err)
		}
	}
	return changed, nil
}

func replaceInLine(line string) (string, error) {
	matches := imageRegex.FindStringSubmatch(line)
	if len(matches) < 6 {
		return line, nil
	}

	prefix := matches[1]
	quote1 := matches[2]
	ref := matches[3]
	quote2 := matches[4]
	suffix := matches[5]

	if isSHA(ref) {
		log.WithField("ref", ref).Debug("already pinned")
		return line, nil
	}

	if strings.Contains(ref, "$") || strings.HasPrefix(ref, "\\") || strings.HasPrefix(ref, ".") {
		return line, nil
	}

	digest, err := getDigest(ref)
	if err != nil {
		return line, err
	}

	repo, tag := parseImage(ref)

	newRef := fmt.Sprintf("%s:%s@%s", repo, tag, digest)

	if idx := strings.Index(suffix, "#"); idx > -1 {
		suffix = suffix[:idx]
	}
	suffix = strings.TrimRight(suffix, " \t")

	newLine := fmt.Sprintf("%s%s%s%s%s # %s", prefix, quote1, newRef, quote2, suffix, tag)

	return newLine, nil
}

func getDigest(ref string) (string, error) {
	repo, tag := parseImage(ref)
	fullRef := repo + ":" + tag

	if v, ok := digestCache[fullRef]; ok {
		return v, nil
	}

	digest, err := crane.Digest(fullRef, crane.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return "", fmt.Errorf("crane: %w", err)
	}

	digestCache[fullRef] = digest
	return digest, nil
}

func parseImage(ref string) (string, string) {
	repo, tag, ok := strings.Cut(ref, ":")
	if !ok {
		tag = "latest"
	}

	if idx := strings.Index(tag, "@"); idx != -1 {
		tag = tag[:idx]
	}
	return repo, tag
}

func isSHA(s string) bool {
	_, digest, found := strings.Cut(s, "@sha256:")
	if !found {
		return false
	}
	return len(digest) == shaLen
}

func isSupported(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	ext := filepath.Ext(base)

	if strings.Contains(base, "dockerfile") {
		return true
	}

	if (ext == ".yml" || ext == ".yaml") &&
		(strings.Contains(base, "docker-compose") || strings.HasPrefix(base, "compose")) {
		return true
	}
	return false
}
