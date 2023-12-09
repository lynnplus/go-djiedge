# Go-DJIEdge

[![Go Reference](https://pkg.go.dev/badge/github.com/lynnplus/go-djiedge.svg)](https://pkg.go.dev/github.com/lynnplus/go-djiedge)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/lynnplus/go-djiedge)
![GitHub tag (with filter)](https://img.shields.io/github/v/tag/lynnplus/go-djiedge)
![](https://img.shields.io/badge/platform-linux-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/lynnplus/go-djiedge)](https://goreportcard.com/report/github.com/lynnplus/go-djiedge)
[![GitHub](https://img.shields.io/github/license/lynnplus/go-djiedge)](https://github.com/lynnplus/go-djiedge/blob/master/LICENSE)

The go-djiedge package provides Go language bindings for the dji-edge-sdk.

The package supports all functions included in dji-edge-sdk and provides a more friendly API usage.

For more specific functions, please visit
DJI Edge-SDK document:
https://developer.dji.com/doc/edge-sdk-tutorial/en/

## Dependency

The package only run on Linux systems, it is recommended to use Ubuntu 22.04.1 LTS

- stdc++ >=11 (recommend c++14 or higher)
- openssl >=1.1.1f
- libssh2 >=1.10.0 ,install by command `sudo apt-get install libssh2-1-dev`

For more software requirements, please visit DJI-Edge-SDK
Document: [environment-prepare](https://developer.dji.com/doc/edge-sdk-tutorial/en/quick-start/environment-prepare.html#software-installation)

## Compiler And Install

### Get Package

Run the command in the project console for Install the go-djiedge package:

`go get github.com/lynnplus/go-djiedge`

### Build

The package does not include the header files and precompiled static libraries of DJI Edge-SDK.
you need to download the SDK manually.

```git clone https://github.com/dji-sdk/Edge-SDK.git```

By default, the package only following cgo compilation instructions are provided,
so you need to manually configure the cgo system environment variables

[default_build.go](default_build.go)  
`#cgo LDFLAGS: -ledgesdk -lcrypto -lssh2`

CGO compilation instructions can be temporarily configured in the console.

For example:

```shell
export CGO_CXXFLAGS="-std=c++14 -I{your Edge-SDK path}/include"
export CGO_LDFLAGS="-L{your Edge-SDK path}/lib/{x86_64 or aarch64}"
```

### Custom Build

If you need to customize the build, such as using different compilation instructions, or the name of `libedgesdk.a`
changes,
this behavior can be disabled by supplying `-tags custom_edge_env` when building/running your application.
when building with this tag you will need to supply the complete CGO environment variables yourself.

For example:

```shell
export CGO_CXXFLAGS="-std=c++14 -I{your Edge-SDK path}/include"
export CGO_LDFLAGS="-L{your Edge-SDK path}/lib/{x86_64 or aarch64} -ledgesdk -lcrypto -lssh2"
```

### Stream Simulate

It is very inconvenient to debug real Dji-Edge devices, and Edge-SDK only supports linux (arch:x86_64/aarch64) operating
system.<br>
the package provides a simple simulation implementation,
supports SDK initialization and pushing real-time camera streams,
currently only supports read h264 file[edge_stream.h264]. <br>

Construction constraints that currently enable the simulation function `//go:build !linux || fake_edge`

The implementation file is in [edge_simulation.go](edge_simulation.go),can add `-tags fake_edge` to open it when
building.

Step Tips:

1. Add `fake_edge` build constraints, build executable file<br>
   `go build -tags fake_edge`
2. Convert mp4 file to h264 format file<br>
    ```shell
    sudo apt-get install ffmpeg
    ffmpeg -i test.mp4 -codec copy -bsf: h264_mp4toannexb -f h264 edge_stream.h264
    ```
3. Copy the generated edge_stream.h264 to the same level directory as the executable file
4. Run