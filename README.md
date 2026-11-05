# Restreamer, modernized

This fork turns the former multi-repository Restreamer build into one maintainable
source tree and one Linux AMD64 container image. Core, UI, the patched FFmpeg build,
deployment examples, and release automation now live together instead of requiring
four repositories to agree on what year it is.

The image is designed for multi-camera 3D-printer monitoring and contains one FFmpeg
build with:

- NVIDIA CUDA, NVDEC, and NVENC;
- Intel VA-API and Quick Sync through oneVPL;
- AMD VA-API and Vulkan Video through Mesa;
- V4L2 devices, RTSP, RTMP, SRT, HLS, and WebRTC delivery through MediaMTX;
- Restreamer's JSON progress and HLS sequence patches.

The container can fall back to software processing when no GPU is exposed. Compiled
support does not invent a working host driver, although GPU packaging has spent years
trying to make that distinction feel unnecessarily mysterious.

## Status

The unified image builds and boots successfully with FFmpeg 9.0.2. Its build gate
verifies the expected hardware APIs, encoders, filters, license, and Restreamer
patches before an image can be published. Physical-device throughput and
glass-to-glass latency must still be measured on the target Linux host and cameras.

Published images are intended to use:

```text
ghcr.io/kaibsora/restreamer:latest
```

Until the first successful GitHub release workflow completes, build the image locally
as described below.

## Quick start

Requirements:

- Linux AMD64 host;
- Docker Engine with Compose v2;
- a compatible host GPU driver for hardware acceleration;
- NVIDIA Container Toolkit only when using an NVIDIA GPU.

Clone the repository and start Restreamer plus the LAN MediaMTX relay:

```sh
git clone https://github.com/kaibsora/restreamer.git
cd restreamer
docker compose -f deploy/compose.yaml up -d
```

Open `http://HOST:8080/ui/`. Persistent configuration and data are stored in named
Docker volumes.

For Intel or AMD hardware acceleration:

```sh
RENDER_GID="$(getent group render | cut -d: -f3)" \
VIDEO_GID="$(getent group video | cut -d: -f3)" \
docker compose -f deploy/compose.yaml -f deploy/compose.dri.yaml up -d
```

For NVIDIA hardware acceleration:

```sh
docker compose -f deploy/compose.yaml -f deploy/compose.nvidia.yaml up -d
```

For USB cameras, copy `deploy/compose.camera-example.yaml`, list only the devices the
container needs, and add the copied file as another Compose override. Do not use
`--privileged` merely because a webcam exists; device permissions are not a séance.

Verify what the running container can actually see:

```sh
docker exec restreamer restreamer-hardware-report
```

See [deploy/README.md](deploy/README.md) for device mapping, WebRTC connectivity,
Crowsnest, Moonraker, and low-latency camera topology.

## Camera and latency model

LAN monitoring and public viewing use separate delivery planes:

1. Restreamer ingests USB cameras, IP-camera RTSP URLs, or RTSP produced by
   Crowsnest's `camera-streamer` mode.
2. A LAN output is copied or encoded with a short GOP and published to MediaMTX.
3. Printer dashboards consume MediaMTX through WebRTC; native clients may use RTSP.
4. Public viewers use Restreamer's HLS player through the authenticated reverse
   proxy, where extra buffering is acceptable.

Prefer stream copy for an existing H.264 camera feed. Every unnecessary decode and
re-encode adds latency, GPU load, and another opportunity for software to demonstrate
creativity. When transcoding is required, use a hardware encoder, no B-frames, a
one-second-or-shorter GOP, and the low-latency input option exposed by the UI.

Sub-second latency is a target to validate on the real network, not a property that
can be certified by compiling an encoder inside Docker Desktop.

## Build locally

Build the hardware FFmpeg base first:

```sh
IMAGE_NAME=restreamer-ffmpeg:hardware components/ffmpeg/Build.sh
```

The helper and GitHub Actions share the same stable-release resolver. Pass a
specific release as the first argument only when reproducing an older build.

Then build the unified application:

```sh
docker buildx build --load \
  --build-arg FFMPEG_IMAGE=restreamer-ffmpeg:hardware \
  --tag restreamer:local .
```

Run the same FFmpeg contract used during releases:

```sh
docker run --rm --entrypoint /usr/local/bin/restreamer-verify-ffmpeg \
  restreamer:local
```

Run the complete Restreamer, RTSP relay, readback, and WebRTC smoke test:

```sh
RESTREAMER_IMAGE=restreamer:local scripts/smoke-compose.sh
```

## Releases

`.github/workflows/build-unified.yml` supports manual and scheduled builds. It:

1. resolves the newest stable FFmpeg source release unless a version is supplied;
2. builds the redistributable hardware FFmpeg image;
3. fails if required hardware support or Restreamer patches disappear;
4. rejects high or critical production dependency vulnerabilities;
5. runs the Core tests and vet checks plus all UI tests;
6. validates the Compose and MediaMTX deployment configuration;
7. pins the unified build to the FFmpeg image digest;
8. boots and probes the application before publishing to GitHub Container Registry;
9. emits provenance and an SBOM for both images.

Pull requests run the quality gates without publishing images. The two inherited
bundle workflows were removed because they built external datarhei components with
obsolete FFmpeg images and attempted to publish to the upstream Docker Hub account.
One repository now has one release pipeline, a concept CI systems traditionally
survive without emotional distress.

An upstream FFmpeg release that no longer accepts the Restreamer patches is expected
to fail loudly. Silently publishing an unpatched binary would be impressively tidy
automation of the wrong result.

## Repository layout

- `components/core`: Go backend and process manager
- `components/ui`: React administration and player UI
- `components/ffmpeg`: FFmpeg source build and Restreamer patches
- `deploy`: Compose, MediaMTX, camera, and GPU deployment assets
- `docs/architecture.md`: system architecture and release policy
- `.github/workflows`: image build and publication automation

See [CONTRIBUTING.md](CONTRIBUTING.md) for development and upstream synchronization.

## License and provenance

The repository retains the original Apache-2.0 project license and component
history. FFmpeg is built without nonfree components and reports a redistributable
GPLv3-or-later configuration. See [LICENSE](LICENSE) and the component license files
for details.
