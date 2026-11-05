# Core

![dsdsds](https://github.com/datarhei/misc/blob/main/img/media-core.png?raw=true)

[![License: Apache2](https://img.shields.io/badge/License-Apache%202.0-brightgreen.svg)](<[https://opensource.org/licenses/MI](https://www.apache.org/licenses/LICENSE-2.0)>)
[![CodeQL](https://github.com/datarhei/core/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/datarhei/core/actions/workflows/codeql-analysis.yml)
[![tests](https://github.com/datarhei/core/actions/workflows/go-tests.yml/badge.svg)](https://github.com/datarhei/core/actions/workflows/go-tests.yml)
[![codecov](https://codecov.io/gh/datarhei/core/branch/main/graph/badge.svg?token=90YMPZRAFK)](https://codecov.io/gh/datarhei/core)
[![Go Report Card](https://goreportcard.com/badge/github.com/datarhei/core)](https://goreportcard.com/report/github.com/datarhei/core)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/datarhei/core)](https://pkg.go.dev/github.com/datarhei/core)
[![Gitbook](https://img.shields.io/badge/GitBook-quick%20start-green)](https://docs.datarhei.com/core/guides/beginner)

The datarhei Core is a process management solution for FFmpeg that offers a range of interfaces for media content, including HTTP, RTMP, SRT, and storage options. It is optimized for use in virtual environments such as Docker. It has been implemented in various contexts, from small-scale applications like Restreamer to large-scale, multi-instance frameworks spanning multiple locations, such as dedicated servers, cloud instances, and single-board computers. The datarhei Core stands out from traditional media servers by emphasizing FFmpeg and its capabilities rather than focusing on media conversion.

## Objectives of development

The objectives of development are:

-   Unhindered use of FFmpeg processes
-   Portability of FFmpeg, including management across development and production environments
-   Scalability of FFmpeg-based applications through the ability to offload processes to additional instances
-   Streamlining of media product development by focusing on features and design.

## What issues have been resolved thus far?

### Process management

-   Run multiple processes via API
-   Unrestricted FFmpeg commands in process configuration.
-   Error detection and recovery (e.g., FFmpeg stalls, dumps)
-   Referencing for process chaining (pipelines)
-   Placeholders for storage, RTMP, and SRT usage (automatic credentials management and URL resolution)
-   Logs (access to current stdout/stderr)
-   Log history (configurable log history, e.g., for error analysis)
-   Resource limitation (max. CPU and MEMORY usage per process)
-   Statistics (like FFmpeg progress per input and output, CPU and MEMORY, state, uptime)
-   Input verification (like FFprobe)
-   Metadata (option to store additional information like a title)

### Media delivery

-   Configurable file systems (in-memory, disk-mount, S3)
-   HTTP/S, RTMP/S, and SRT services, including Let's Encrypt
-   Bandwidth and session limiting for HLS/MPEG DASH sessions (protects restreams from congestion)
-   Viewer session API and logging

### Misc

-   HTTP REST and GraphQL API
-   Swagger documentation
-   Metrics incl. Prometheus support (also detects POSIX and cgroups resources)
-   Docker images for fast setup of development environments up to the integration of cloud resources

## Docker images

-   datarhei/core:latest (AMD64, ARM64, ARMv7)
-   datarhei/core:cuda-latest (Nvidia CUDA 11.7.1, AMD64)
-   datarhei/core:rpi-latest (Raspberry Pi / OMX/V4L2-M2M, AMD64/ARMv7)
-   datarhei/core:vaapi-latest (Intel VAAPI, AMD64)

## Quick start

1. Run the Docker image

```sh
docker run --name core -d \
    -e CORE_API_AUTH_USERNAME=admin \
    -e CORE_API_AUTH_PASSWORD=secret \
    -p 8080:8080 \
    -v ${HOME}/core/config:/core/config \
    -v ${HOME}/core/data:/core/data \
    datarhei/core:latest
```

2. Open Swagger
   http://host-ip:8080/api/swagger/index.html

3. Log in with Swagger
   Authorize > Basic authorization > Username: admin, Password: secret

## Documentation

Documentation is available on [docs.datarhei.com/core](https://docs.datarhei.com/core).

-   [Quick start](https://docs.datarhei.com/core/guides/beginner)
-   [Installation](https://docs.datarhei.com/core/installation)
-   [Configuration](https://docs.datarhei.com/core/configuration)
-   [Coding](https://docs.datarhei.com/core/development/coding)

## License

datarhei/core is licensed under the Apache License 2.0
