# piñker 🧷

[![Build Status](https://img.shields.io/github/actions/workflow/status/honeok/pinker/build.yaml?branch=master&logo=github)](https://github.com/honeok/pinker)
[![Go Report](https://goreportcard.com/badge/github.com/honeok/pinker)](https://goreportcard.com/report/github.com/honeok/pinker)
[![GitHub Release](https://img.shields.io/github/release/honeok/pinker.svg?logo=github)](https://github.com/honeok/pinker/releases/latest)
[![GitHub License](https://img.shields.io/github/license/honeok/pinker.svg?logo=github)](https://github.com/honeok/pinker)

Secure your container supply chain by automatically pinning Docker images to immutable SHA256 digests.

> [!TIP]
> Tags are mutable. Digests are forever.

<img width="1930" height="994" src="https://github.com/user-attachments/assets/d840ed52-2013-4625-8d1a-2f8abbd2a5f1" />

## Why Pinker?

Mutable tags like `postgres:14` or `node:latest` can change at any time. Pinker resolves these tags to their exact SHA256 digest, ensuring your builds are:

- **Reproducible**: Everyone builds the exact same image, every time.
- **Secure**: Protects against tag hijacking or unexpected upstream changes.

## Features

- Zero-Config Authentication Seamlessly reuses your local Docker credentials. If `docker login registry.example.com` works in your terminal, `pinker` works automatically. No need to manage extra secrets.
- Multi-Registry Support Built on top of standard OCI libraries with native support for **AWS ECR**, **Google Artifact Registry**, **Azure ACR**, **Harbor**, and any OCI-compliant private registry.
- **Immutable by Design** Converts uncertain tags (e.g., `latest`) into deterministic SHA256 digests, locking your infrastructure to a specific state.
- **Broad Compatibility** Works out-of-the-box with `Dockerfile` and `docker-compose.yml` (supports `.yaml`, `compose.yml`, and other standard naming conventions).

## Install

Binaries (Recommended)

Download the pre-compiled binaries for your platform from the [Releases page][1].

Go Install

If you have Go installed, you can build and install the latest version from source:

```shell
go install github.com/honeok/pinker@latest
```

## Usage

```shell
# Pin everything in the current directory (recursive)
$ pinker

# Pin a specific directory
$ pinker ./deploy
```

## Acknowledgements

This project is heavily inspired by [pinata][2] by [caarlos0][3]. The code structure, design philosophy, and CLI experience are directly adapted from his work to bring the same security and reproducibility standards to the Docker ecosystem.

[1]: https://github.com/honeok/pinker/releases
[2]: https://github.com/caarlos0/pinata
[3]: https://github.com/caarlos0
