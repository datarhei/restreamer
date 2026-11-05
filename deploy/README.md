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

## LAN RTSP and WebRTC

The deployment includes MediaMTX as the LAN-facing real-time origin. In Restreamer,
add an RTSP publication for each camera and use a unique destination such as:

```text
rtsp://mediamtx:8554/printer-one
rtsp://mediamtx:8554/printer-two
```

Prefer stream copy when the camera already emits browser-compatible H.264. If it
must be encoded, select a low-latency hardware encoder, a one-second GOP, no B-frames,
and RTSP over TCP for a reliable wired LAN or UDP where minimum delay matters more
than packet-loss recovery.

Other LAN clients can read each stream directly:

```text
rtsp://RESTREAMER_LAN_IP:8554/printer-one
http://RESTREAMER_LAN_IP:8889/printer-one
```

For Mainsail/Moonraker, configure the WebRTC endpoint rather than handing a browser
an RTSP URL it cannot play:

```ini
[webcam printer_one]
service: webrtc-mediamtx
stream_url: http://RESTREAMER_LAN_IP:8889/printer-one
snapshot_url:
target_fps: 25
```

Moonraker identifies `webrtc-mediamtx` as a supported Mainsail webcam service. The
printer does not need Crowsnest to ingest this remote stream; Crowsnest can continue
handling locally attached cameras while Moonraker/Mainsail displays either source.

The RTSP, RTP, RTCP, WebRTC HTTP, and WebRTC ICE ports are intended for the trusted
LAN. Do not publish them through the public reverse proxy without authentication and
firewall rules. Public viewers should continue using Restreamer's authenticated HLS
player path, where several seconds of buffering are useful rather than tragic.
