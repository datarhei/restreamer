# FFmpeg Base

FFmpeg base image for [datarhei/core](https://github.com/datarhei/core).

[![alpine](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_alpine.yaml/badge.svg)](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_alpine.yaml)
[![alpine-rpi](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_alpine-rpi.yaml/badge.svg)](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_alpine-rpi.yaml)
[![base:ubuntu-ffmpeg-vvapi](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_ubuntu-vaapi.yaml/badge.svg)](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_ubuntu-vaapi.yaml)
[![ubuntu-cuda](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_ubuntu-cuda.yaml/badge.svg)](https://github.com/datarhei/ffmpeg/actions/workflows/build_base_ubuntu-cuda.yaml)

Branch: 6.1

## Config:

```sh
--enable-libv4l2
--enable-libfreetype
--enable-alsa
--enable-libsrt
--enable-libx264
--enable-libx265
--enable-libvpx
--enable-libmp3lame
--enable-libopus
--enable-libvorbis
```

_Additional informations can be found in the Dockerfiles._

## Patches ([contrib](contrib/)):

- JSON-Stats (expands progress data per file in json format)
- HLS Bitrate (calculates bitrate estimate for HLS master playlist)

## Images and Plattforms:

| Dockerimage                                        | OS           | Plattform                                | GPU                                         |
| -------------------------------------------------- | ------------ | ---------------------------------------- | ------------------------------------------- |
| docker.io/datarhei/base:alpine-ffmpeg-latest       | Alpine 3.16  | linux/amd64, linux/arm64, linux/arm/v7   | -                                           |
| docker.io/datarhei/base:alpine-ffmpeg-rpi-latest   | Alpine 3.16  | Raspberry Pi (linux/arm/v7, linux/arm64) | MMAL/OMX/V4L2-M2M (32bit), V4L2-M2M (64bit) |
| docker.io/datarhei/base:ubuntu-ffmpeg-cuda-latest  | Ubuntu 20.04 | linux/amd64                              | Nvidia Cuda                                 |
| docker.io/datarhei/base:ubuntu-ffmpeg-vaapi-latest | Ubuntu 20.04 | linux/amd64                              | Intel VAAPI                                 |

More tags: https://hub.docker.com/repository/docker/datarhei/base/general

## Build & test

```sh
$ git clone github.com/datarhei/ffmpeg
$ ./Build.sh {arg}
```

Args:

- default (alpine-ffmpeg-latest)
- rpi (alpine-ffmpeg-rpi-latest)
- cuda (ubuntu-ffmpeg-cuda-latest)
- vaapi (ubuntu-ffmpeg-vaapi-latest)

## Known problems:

The libraries are currently not compiled due to errors caused by Docker virtualisation.

## Feature requests:

Please create an issue with your use case and all the requirements.

## Licence

LGPL-licensed with optional components licensed under GPL. Please refer to the LICENSE file for detailed information.
