# piñker

[![Build Status](https://img.shields.io/github/actions/workflow/status/honeok/pinker/build.yaml?branch=master&logo=github)](https://github.com/honeok/pinker)
[![GitHub Release](https://img.shields.io/github/release/honeok/pinker.svg?logo=github)](https://github.com/honeok/pinker/releases/latest)
[![GitHub License](https://img.shields.io/github/license/honeok/pinker.svg?logo=github)](https://github.com/honeok/pinker)

Make your Docker infrastructure more secure and reproducible by pinning images to their SHA256 digests.

> [!Note]
>
> Just like pinata for GitHub Actions, but for Dockerfiles and Compose files.

<img width="1930" height="994" src="https://github.com/user-attachments/assets/d840ed52-2013-4625-8d1a-2f8abbd2a5f1" />

## Features

- Zero-Config Authentication No need to hardcode secrets or mount volumes. `pinker` seamlessly reuses your local credentials. If you have run `docker login registry.example.com` in your terminal, `pinker` automatically has access.
- Multi-Cloud Support Built on top of standard OCI libraries, offering native support for **AWS ECR**, **Google Artifact Registry**, **Azure ACR**, **Harbor**, and any other private registry out of the box.
- Results-Oriented It does one thing and does it well: converting mutable, uncertain tags (e.g., `latest`) into immutable, deterministic SHA256 digests.
- **Broad Compatibility** Works with both `Dockerfile` and `docker-compose.yml` (including variants like `.yaml`, `compose.yml`, etc.).

## Install

```shell
# Go:
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

This project is heavily inspired by [pinata][1] by [caarlos0][2].
The code structure, design philosophy, and CLI experience are directly adapted from his work to bring the same security and reproducibility standards to the Docker ecosystem.

Thank you, Carlos, for the inspiration! ❤️

[1]: https://github.com/caarlos0/pinata
[2]: https://github.com/caarlos0
