# Container deployment

Start the CPU-only container:

```sh
docker compose -f deploy/compose.yaml up -d
```

For Intel or AMD acceleration, expose the DRM render nodes and set the host's
actual `video` and `render` group IDs:

```sh
RENDER_GID="$(getent group render | cut -d: -f3)" \
VIDEO_GID="$(getent group video | cut -d: -f3)" \
docker compose -f deploy/compose.yaml -f deploy/compose.dri.yaml up -d
```

For NVIDIA acceleration, install NVIDIA Container Toolkit on the host and use:

```sh
docker compose -f deploy/compose.yaml -f deploy/compose.nvidia.yaml up -d
```

Copy `compose.camera-example.yaml`, enumerate every required `/dev/videoN` device,
and include the copy as another Compose override. Device numbers may change after a
reboot; persistent udev symlinks are recommended for more than one USB camera.

After startup, inspect what the container can actually use:

```sh
docker exec restreamer restreamer-hardware-report
```

An encoder appearing in `ffmpeg -encoders` proves it was compiled. Successful
`vainfo`, `vulkaninfo`, or `nvidia-smi` output proves that the corresponding host
runtime reached the container. Those are separate facts, despite decades of GPU
software trying to make them feel like one exciting mystery.
