#!/bin/bash

set -au

function build_default_native() {
  export OS_NAME=alpine
  export OS_VERSION=3.23
  export FFMPEG_VERSION=8.1.2

  docker build \
    --progress=plain \
    --build-arg BUILD_IMAGE=$OS_NAME:$OS_VERSION \
    --build-arg FFMPEG_VERSION=$FFMPEG_VERSION \
    -f Dockerfile.alpine \
    -t datarhei/base:ffmpeg${FFMPEG_VERSION}-${OS_NAME}${OS_VERSION} .
}

function build_default() {
  export OS_NAME=alpine
  export OS_VERSION=3.23
  export FFMPEG_VERSION=8.1.2

  docker buildx build \
    --load \
    --progress=plain \
    --build-arg BUILD_IMAGE=$OS_NAME:$OS_VERSION \
    --build-arg FFMPEG_VERSION=$FFMPEG_VERSION \
    --platform linux/amd64 \
    -f Dockerfile.alpine \
    -t datarhei/base:ffmpeg${FFMPEG_VERSION}-${OS_NAME}${OS_VERSION} .
}

function build_rpi() {
  export OS_NAME=alpine
  export OS_VERSION=3.21
  export FFMPEG_VERSION=8.1.2

  docker build \
    --progress=plain \
    --build-arg BUILD_IMAGE=$OS_NAME:$OS_VERSION \
    --build-arg FFMPEG_VERSION=$FFMPEG_VERSION \
    -f Dockerfile.alpine.rpi \
    -t datarhei/base:ffmpeg${FFMPEG_VERSION}-rpi-${OS_NAME}${OS_VERSION} .
}

function build_cuda12() {
  export OS_NAME=ubuntu
  export OS_VERSION=24.04
  export FFMPEG_VERSION=8.1.2
  export FFNVCODEC_VERSION=12.2.72.0

  docker build \
    --progress=plain \
    --build-arg BUILD_IMAGE=nvidia/cuda:$CUDA_VERSION-devel-ubuntu$OS_VERSION \
    --build-arg DEPLOY_IMAGE=nvidia/cuda:$CUDA_VERSION-runtime-ubuntu$OS_VERSION \
    --build-arg FFNVCODEC_VERSION=$FFNVCODEC_VERSION \
    --build-arg FFMPEG_VERSION=$FFMPEG_VERSION \
    -f Dockerfile.ubuntu.cuda12 \
    -t datarhei/base:ffmpeg${FFMPEG_VERSION}-cuda-ubuntu$OS_VERSION-cuda${CUDA_VERSION} .
}

function build_cuda13() {
  export OS_NAME=ubuntu
  export OS_VERSION=26.04
  export FFMPEG_VERSION=8.1.2
  export FFNVCODEC_VERSION=13.0.19.0

  docker build \
    --progress=plain \
    --build-arg BUILD_IMAGE=nvidia/cuda:$CUDA_VERSION-devel-ubuntu$OS_VERSION \
    --build-arg DEPLOY_IMAGE=nvidia/cuda:$CUDA_VERSION-runtime-ubuntu$OS_VERSION \
    --build-arg FFNVCODEC_VERSION=$FFNVCODEC_VERSION \
    --build-arg FFMPEG_VERSION=$FFMPEG_VERSION \
    -f Dockerfile.ubuntu.cuda13 \
    -t datarhei/base:ffmpeg${FFMPEG_VERSION}-cuda-ubuntu$OS_VERSION-cuda${CUDA_VERSION} .
}

function build_vaapi() {
  export OS_NAME=ubuntu
  export OS_VERSION=26.04
  export FFMPEG_VERSION=8.1.2

  docker buildx build \
    --load \
    --progress=plain \
    --build-arg BUILD_IMAGE=$OS_NAME:$OS_VERSION \
    --build-arg DEPLOY_IMAGE=$OS_NAME:$OS_VERSION \
    --build-arg FFMPEG_VERSION=$FFMPEG_VERSION \
    --platform linux/amd64 \
    -f Dockerfile.ubuntu.vaapi \
    -t datarhei/base:ffmpeg${FFMPEG_VERSION}-vaapi-${OS_NAME}${OS_VERSION} .
}

main() {
  if [[ $# == 0 ]]; then
    echo "Options available: default, default_native, rpi, cuda12, vaapi"
    exit 0
  else
    if [[ $1 == "default" ]]; then
      build_default
    elif [[ $1 == "default_native" ]]; then
      build_default_native
    elif [[ $1 == "rpi" ]]; then
      build_rpi
    elif [[ $1 == "cuda12" ]]; then
      build_cuda12
    elif [[ $1 == "cuda13" ]]; then
      build_cuda13
    elif [[ $1 == "vaapi" ]]; then
      build_vaapi
    fi
  fi
}

main $@

exit 0
