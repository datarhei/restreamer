#!/usr/bin/env bash
set -euo pipefail

project_name="${COMPOSE_PROJECT_NAME:-restreamer-integration-smoke}"
restreamer_image="${RESTREAMER_IMAGE:-restreamer-smoke:local}"
compose=(docker compose --project-name "${project_name}" --file deploy/compose.yaml)

if docker container inspect restreamer >/dev/null 2>&1 ||
    docker container inspect restreamer-mediamtx >/dev/null 2>&1; then
    printf '%s\n' 'Refusing to replace an existing restreamer smoke-test container.' >&2
    exit 1
fi

cleanup() {
    status=$?
    trap - EXIT
    if [[ ${status} -ne 0 ]]; then
        "${compose[@]}" logs --no-color --tail 100 || true
    fi
    "${compose[@]}" down --volumes --remove-orphans >/dev/null 2>&1 || true
    exit "${status}"
}
trap cleanup EXIT

export RESTREAMER_IMAGE="${restreamer_image}"
export RESTREAMER_WEBRTC_HOSTS="${RESTREAMER_WEBRTC_HOSTS:-127.0.0.1}"

"${compose[@]}" up --detach

ready=false
for _ in $(seq 1 30); do
    if curl --fail --silent http://127.0.0.1:8080/api >/dev/null; then
        ready=true
        break
    fi
    sleep 1
done
[[ "${ready}" == "true" ]]

for path in / /ui/ /api /api/v3/skills; do
    curl --fail --silent --show-error "http://127.0.0.1:8080${path}" >/dev/null
done

docker exec restreamer restreamer-verify-ffmpeg

docker exec --detach restreamer \
    ffmpeg -hide_banner -loglevel error -re \
    -f lavfi -i testsrc2=size=320x180:rate=15 \
    -an -c:v libx264 -preset ultrafast -tune zerolatency -g 15 -bf 0 \
    -f rtsp -rtsp_transport tcp rtsp://mediamtx:8554/ci-smoke

stream_info=''
for _ in $(seq 1 15); do
    if stream_info=$(docker exec restreamer \
        ffprobe -hide_banner -v error -rtsp_transport tcp \
        -select_streams v:0 \
        -show_entries stream=codec_name,width,height,r_frame_rate \
        -of default=noprint_wrappers=1 \
        rtsp://mediamtx:8554/ci-smoke 2>/dev/null); then
        break
    fi
    sleep 1
done

grep -Fx 'codec_name=h264' <<<"${stream_info}" >/dev/null
grep -Fx 'width=320' <<<"${stream_info}" >/dev/null
grep -Fx 'height=180' <<<"${stream_info}" >/dev/null
grep -Fx 'r_frame_rate=15/1' <<<"${stream_info}" >/dev/null

curl --fail --silent --show-error \
    http://127.0.0.1:8889/ci-smoke >/dev/null

printf '%s\n' 'Unified Restreamer + RTSP + MediaMTX WebRTC smoke test passed.'
