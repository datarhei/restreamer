#!/bin/sh
set -eu

verify=false
case "${1:-}" in
    "") ;;
    --verify) verify=true ;;
    *)
        echo "Usage: restreamer-hardware-report [--verify]" >&2
        exit 64
        ;;
esac

workdir="$(mktemp -d)"
trap 'rm -rf "${workdir}"' EXIT INT TERM

render_device=''
va_info=''
vulkan_info=''
nvidia_info=''

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
    for device in /dev/dri/renderD*; do
        [ -e "${device}" ] || continue
        if vainfo --display drm --device "${device}" >"${workdir}/vainfo" 2>&1; then
            render_device="${device}"
            va_info="$(cat "${workdir}/vainfo")"
            printf 'Render device: %s\n%s\n' "${render_device}" "${va_info}"
            break
        fi
    done
    if [ -z "${render_device}" ]; then
        echo "No usable VA-API render device detected."
    fi
fi

if command -v vulkaninfo >/dev/null 2>&1; then
    echo
    echo "== Vulkan runtime =="
    if vulkaninfo --summary >"${workdir}/vulkaninfo" 2>&1; then
        vulkan_info="$(cat "${workdir}/vulkaninfo")"
        printf '%s\n' "${vulkan_info}"
    else
        cat "${workdir}/vulkaninfo"
    fi
fi

if command -v nvidia-smi >/dev/null 2>&1; then
    echo
    echo "== NVIDIA runtime =="
    if nvidia-smi >"${workdir}/nvidia-smi" 2>&1; then
        nvidia_info="$(cat "${workdir}/nvidia-smi")"
        printf '%s\n' "${nvidia_info}"
    else
        cat "${workdir}/nvidia-smi"
    fi
fi

if [ "${verify}" != true ]; then
    exit 0
fi

echo
echo "== Active hardware encode verification =="

runtime_count=0
failure_count=0

verify_encoder() {
    label=$1
    shift
    runtime_count=$((runtime_count + 1))

    if "$@" >"${workdir}/encode" 2>&1; then
        printf 'PASS: %s encoded a test frame.\n' "${label}"
    else
        printf 'FAIL: %s could not encode a test frame.\n' "${label}" >&2
        sed 's/^/  /' "${workdir}/encode" >&2
        failure_count=$((failure_count + 1))
    fi
}

if [ -n "${nvidia_info}" ]; then
    verify_encoder "NVIDIA NVENC" \
        ffmpeg -hide_banner -loglevel error \
        -f lavfi -i color=size=128x72:rate=1 \
        -frames:v 1 -c:v h264_nvenc -f null -
fi

if [ -n "${render_device}" ]; then
    verify_encoder "VA-API (${render_device})" \
        ffmpeg -hide_banner -loglevel error \
        -init_hw_device "vaapi=va:${render_device}" -filter_hw_device va \
        -f lavfi -i color=size=128x72:rate=1 \
        -vf format=nv12,hwupload -frames:v 1 -c:v h264_vaapi -f null -

    if printf '%s\n' "${va_info}" | grep -Eqi 'Intel|iHD|i965'; then
        verify_encoder "Intel Quick Sync (${render_device})" \
            ffmpeg -hide_banner -loglevel error \
            -init_hw_device "vaapi=va:${render_device}" \
            -init_hw_device qsv=qsv@va -filter_hw_device qsv \
            -f lavfi -i color=size=128x72:rate=1 \
            -vf format=nv12,hwupload=extra_hw_frames=64 \
            -frames:v 1 -c:v h264_qsv -f null -
    fi
fi

if printf '%s\n' "${vulkan_info}" \
    | grep -Eq 'PHYSICAL_DEVICE_TYPE_(DISCRETE|INTEGRATED)_GPU'; then
    vulkan_device="${RESTREAMER_VULKAN_DEVICE:-0}"
    verify_encoder "Vulkan Video (device ${vulkan_device})" \
        ffmpeg -hide_banner -loglevel error \
        -init_hw_device "vulkan=vk:${vulkan_device}" -filter_hw_device vk \
        -f lavfi -i color=size=128x72:rate=1 \
        -vf format=nv12,hwupload -frames:v 1 -c:v h264_vulkan -f null -
fi

if [ "${runtime_count}" -eq 0 ]; then
    echo "FAIL: no usable hardware-encoding runtime was detected." >&2
    exit 2
fi

if [ "${failure_count}" -ne 0 ]; then
    printf 'Hardware verification failed for %s runtime(s).\n' \
        "${failure_count}" >&2
    exit 1
fi

printf 'Hardware verification passed for %s runtime(s).\n' "${runtime_count}"
