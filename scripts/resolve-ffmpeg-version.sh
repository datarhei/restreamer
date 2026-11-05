#!/usr/bin/env bash

set -euo pipefail

requested_version="${1:-${FFMPEG_VERSION:-}}"
if [[ -n "${requested_version}" ]]; then
	printf '%s\n' "${requested_version#v}"
	exit 0
fi

release_index="$(curl --fail --silent --show-error https://ffmpeg.org/releases/)"
version="$(
	printf '%s' "${release_index}" \
		| grep -oE 'ffmpeg-[0-9]+\.[0-9]+(\.[0-9]+)?\.tar\.xz' \
		| sed -E 's/^ffmpeg-|\.tar\.xz$//g' \
		| awk -F. '
			{
				major = $1 + 0
				minor = $2 + 0
				patch = $3 + 0
				if (!found || major > best_major ||
					(major == best_major && minor > best_minor) ||
					(major == best_major && minor == best_minor && patch > best_patch)) {
					best = $0
					best_major = major
					best_minor = minor
					best_patch = patch
					found = 1
				}
			}
			END { if (found) print best }
		'
)"

if [[ -z "${version}" ]]; then
	echo 'Unable to resolve the latest stable FFmpeg release.' >&2
	exit 1
fi

printf '%s\n' "${version}"
