#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPOSITORY_ROOT="$(cd -- "${SCRIPT_DIR}/../.." && pwd)"

FFMPEG_VERSION="${1:-${FFMPEG_VERSION:-}}"
if [[ -z "${FFMPEG_VERSION}" ]]; then
	FFMPEG_VERSION="$("${REPOSITORY_ROOT}/scripts/resolve-ffmpeg-version.sh")"
fi
FFMPEG_VERSION="${FFMPEG_VERSION#v}"

IMAGE_NAME="${IMAGE_NAME:-restreamer-ffmpeg:hardware-${FFMPEG_VERSION}}"
BUILD_IMAGE="${BUILD_IMAGE:-nvidia/cuda:12.9.1-devel-ubuntu24.04}"
DEPLOY_IMAGE="${DEPLOY_IMAGE:-nvidia/cuda:12.9.1-runtime-ubuntu24.04}"
FFNVCODEC_VERSION="${FFNVCODEC_VERSION:-12.2.72.0}"
VULKAN_HEADERS_VERSION="${VULKAN_HEADERS_VERSION:-1.4.362}"
BUILD_DATE="${BUILD_DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"
BUILD_COMMIT="${BUILD_COMMIT:-$(git -C "${REPOSITORY_ROOT}" rev-parse --short=12 HEAD)}"

echo "Building ${IMAGE_NAME} with FFmpeg ${FFMPEG_VERSION}"

docker buildx build \
	--load \
	--file "${SCRIPT_DIR}/Dockerfile.ubuntu.hardware" \
	--build-arg "BUILD_IMAGE=${BUILD_IMAGE}" \
	--build-arg "DEPLOY_IMAGE=${DEPLOY_IMAGE}" \
	--build-arg "BUILD_DATE=${BUILD_DATE}" \
	--build-arg "BUILD_COMMIT=${BUILD_COMMIT}" \
	--build-arg "FFNVCODEC_VERSION=${FFNVCODEC_VERSION}" \
	--build-arg "VULKAN_HEADERS_VERSION=${VULKAN_HEADERS_VERSION}" \
	--build-arg "FFMPEG_VERSION=${FFMPEG_VERSION}" \
	--tag "${IMAGE_NAME}" \
	"${SCRIPT_DIR}"
