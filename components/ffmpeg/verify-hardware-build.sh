#!/bin/sh
set -eu

fail() {
    printf 'FFmpeg hardware verification failed: %s\n' "$1" >&2
    exit 1
}

require_fixed_line() {
    output=$1
    expected=$2
    printf '%s\n' "$output" | grep -F -x -- "$expected" >/dev/null || fail "missing ${expected}"
}

require_pattern() {
    output=$1
    expected=$2
    printf '%s\n' "$output" | grep -E -- "$expected" >/dev/null || fail "missing pattern ${expected}"
}

ffmpeg_version=$(ffmpeg -hide_banner -version 2>&1)
ffprobe_version=$(ffprobe -hide_banner -version 2>&1)
license=$(ffmpeg -hide_banner -L 2>&1)
hwaccels=$(ffmpeg -hide_banner -hwaccels 2>&1)
encoders=$(ffmpeg -hide_banner -encoders 2>&1)
filters=$(ffmpeg -hide_banner -filters 2>&1)
full_help=$(ffmpeg -hide_banner -h full 2>&1)
hls_help=$(ffmpeg -hide_banner -h muxer=hls 2>&1)

require_pattern "$ffmpeg_version" '^ffmpeg version [0-9]'
require_pattern "$ffprobe_version" '^ffprobe version [0-9]'

if printf '%s\n' "$license" | grep -F -i -- 'unredistributable' >/dev/null; then
    fail 'binary is marked unredistributable'
fi
require_pattern "$license" 'GNU General Public License'

for accelerator in cuda vaapi qsv drm vulkan; do
    require_fixed_line "$hwaccels" "$accelerator"
done

for codec in av1 h264 hevc; do
    for backend in nvenc qsv vaapi vulkan; do
        require_pattern "$encoders" "[[:space:]]${codec}_${backend}[[:space:]]"
    done
done
require_pattern "$encoders" '[[:space:]]aac[[:space:]]+AAC'

for filter in overlay_vulkan scale_vulkan transpose_vulkan; do
    require_pattern "$filters" "[[:space:]]${filter}[[:space:]]"
done

require_pattern "$full_help" '^-jsonstats[[:space:]]'
require_pattern "$hls_help" '[[:space:]]adjust_sequence[[:space:]]'

printf '%s\n' "FFmpeg hardware verification passed: $(printf '%s\n' "$ffmpeg_version" | sed -n '1p')"
