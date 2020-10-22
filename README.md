# alpimager

Build custom Alpine Linux images with Docker.

![Go CI](https://github.com/pojntfx/alpimager/workflows/Go%20CI/badge.svg)

[![introduction video](https://img.youtube.com/vi/pxaqts3eHMM/maxresdefault.jpg)](https://youtu.be/pxaqts3eHMM)

## Overview

This projects builds custom Alpine Linux images with the [alpine-make-vm-image utility](https://github.com/alpinelinux/alpine-make-vm-image), but it uses Docker and a simplified interface so that it can run easily on systems other than Alpine Linux that support Docker, such as for example macOS.

## Installation

Linux and macOS binaries are available on [GitHub Releases](https://github.com/pojntfx/alpimager/releases). Windows is no longer supported since WSL2, the new Docker backend in Windows, [uses a Kernel without `nbd` support](https://github.com/microsoft/WSL/issues/5968); if just want a quick Alpine VM on Windows, I recommend using [Alpine WSL](https://www.microsoft.com/en-us/p/alpine-wsl/9p804crf0395) instead until Microsoft resolve the issue.

## Usage

See [testdata](testdata) for example files.

```bash
% alpimager -help
Usage of alpimager:
  -debug
        Enable debugging output
  -output string
        Output image file (default "alpine.qcow2")
  -packages string
        Package list file (default "packages.txt")
  -repositories string
        Repository list file (default "repositories.txt")
  -script string
        Setup script file (default "setup.sh")
```

## License

alpimager (c) 2020 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
