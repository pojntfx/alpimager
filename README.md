# alpimager

Build custom Alpine Linux images with Docker.

![Go CI](https://github.com/pojntfx/alpimager/workflows/Go%20CI/badge.svg)

[![introduction video](https://img.youtube.com/vi/pxaqts3eHMM/maxresdefault.jpg)](https://youtu.be/pxaqts3eHMM)

## Overview

This projects builds custom Alpine Linux images with the [alpine-make-vm-image utility](https://github.com/alpinelinux/alpine-make-vm-image), but it uses Docker and a simplified interface so that it can run easily on systems other than Alpine Linux that support Docker, such as for example macOS.

## Installation

Linux, macOS and Windows binaries are available on [GitHub Releases](https://github.com/pojntfx/alpimager/releases).

## Usage

See [testdata](testdata) for example files.

```bash
% alpimager -help
Usage of alpimager:
  -output string
        Output image file (default "alpine.qcow2")
  -packages string
        Package list file (default "packages.txt")
  -repositories string
        Repository list file (default "repositories.txt")
  -script string
        Setup script file (default "setup.sh")
```

If you are for some reason still using Windows, please:

- Use absolute file paths in the parameters (i.e. `C:\Users\pojntfx\Downloads\packages.txt`)
- Make sure that the files (i.e. `packages.txt`) use proper `LF` line endings

## License

alpimager (c) 2020 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
