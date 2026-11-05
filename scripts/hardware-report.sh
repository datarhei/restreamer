#!/bin/sh
set -eu

echo "== FFmpeg =="
ffmpeg -hide_banner -version | head -n 3

echo
echo "== Compiled hardware APIs =="
ffmpeg -hide_banner -hwaccels

echo
echo "== Hardware video encoders =="
ffmpeg -hide_banner -encoders 2>/dev/null \
    | grep -E '_(nvenc|vaapi|qsv|vulkan|v4l2m2m)([[:space:]]|$)' \
    || echo "No matching hardware encoders were compiled."

echo
echo "== Video devices =="
for device in /dev/video* /dev/dri/renderD*; do
    if [ -e "${device}" ]; then
        ls -l "${device}"
    fi
done

if command -v vainfo >/dev/null 2>&1; then
    echo
    echo "== VA-API runtime =="
    vainfo --display drm 2>&1 || true
fi

if command -v vulkaninfo >/dev/null 2>&1; then
    echo
    echo "== Vulkan runtime =="
    vulkaninfo --summary 2>&1 || true
fi

if command -v nvidia-smi >/dev/null 2>&1; then
    echo
    echo "== NVIDIA runtime =="
    nvidia-smi 2>&1 || true
fi
