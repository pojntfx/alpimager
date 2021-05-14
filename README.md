# alpimager

Build custom Alpine Linux images with Docker.

![Go CI](https://github.com/pojntfx/alpimager/workflows/Go%20CI/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/alpimager.svg)](https://pkg.go.dev/github.com/pojntfx/alpimager)

[![Introduction video](https://img.youtube.com/vi/pxaqts3eHMM/maxresdefault.jpg)](https://youtu.be/pxaqts3eHMM)

## Overview

This projects builds custom Alpine Linux images with the [alpine-make-vm-image utility](https://github.com/alpinelinux/alpine-make-vm-image), but it uses Docker and a simplified interface so that it can run easily on systems other than Alpine Linux that support Docker, such as for example macOS or other Linux distros.

## Installation

Binaries are built weekly and uploaded to [GitHub releases](https://github.com/pojntfx/alpimager/releases).

On Linux, you can install them like so:

```shell
$ curl -L -o /tmp/alpimager https://github.com/pojntfx/alpimager/releases/download/latest/alpimager.linux-$(uname -m)
$ sudo install /tmp/alpimager /usr/local/bin
```

On macOS, you can use the following to install:

```shell
$ curl -L -o /tmp/alpimager https://github.com/pojntfx/alpimager/releases/download/latest/alpimager.darwin-$(uname -m)
$ sudo install /tmp/alpimager /usr/local/bin
```

On Windows, the following should work (using PowerShell as administrator):

```shell
PS> Invoke-WebRequest https://github.com/pojntfx/alpimager/releases/download/latest/alpimager.windows-x86_64.exe -OutFile \Windows\System32\alpimager.exe
```

Note that the **Windows builds are broken** due to a change in how WSL works; it is no longer supported since WSL2, the new Docker backend in Windows, [uses a Kernel without `nbd` support](https://github.com/microsoft/WSL/issues/5968). If you just want a quick Alpine VM on Windows, I recommend using [Alpine WSL](https://alpimagerw.microsoft.com/en-us/p/alpine-wsl/9p804crf0395) instead until Microsoft resolve the issue.

## Usage

See [testdata](./testdata) for example files.

```bash
$ alpimager --help
Usage of alpimager:
  -maximumDiskSize string
        Maximum disk size (default "20G")
  -output string
        Output image file (default "alpine.qcow2")
  -packages string
        Package list file (default "packages.txt")
  -repositories string
        Repository list file (default "repositories.txt")
  -script string
        Setup script file (default "setup.sh")
  -verbose
        Enable verbose logging
```

## Troubleshooting

- If you get an error like `ERROR: No available nbd device found!`, run `sudo modprobe nbd` on the Docker host.
- Running alpimager with `-verbose` might give you more debugging output.

## License

alpimager (c) 2021 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
