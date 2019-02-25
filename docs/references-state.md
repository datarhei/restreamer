---
title: State
---

The Restreamer is storing its state in a JSON file `v1.json` in the mounted volume `/restreamer/db` (the one from the docker command line, e.g. `docker run ... -v /mnt/restreamer/db:/restreamer/db`).

In order to prefill the RTSP and RTMP URLs, you can provide a pre-populated state file with the name `v1.json` in the mounted directory. The format of the file is:

```json
{
    "addresses": {
        "optionalOutputAddress": "rtmp://...",
        "srcAddress": "rtsp://..."
    },
    "options": {
        "rtspTcp": true
    },
    "states": {
        "repeatToLocalNginx": {
            "type": "connected"
        },
        "repeatToOptionalOutput": {
            "type": "connected"
        }
    },
    "userActions": {
        "repeatToLocalNginx": "start",
        "repeatToOptionalOutput": "start"
    }
}
```
This will pull a stream from a RTSP source and repeat it to a RTMP target. Also, it will make the Restreamer reconnecting to the source and target immediately starting up.

The different sections are explained here in more detail.
```json
"addresses": {
    "optionalOutputAddress": "rtmp://...",
    "srcAddress": "rtsp://..."
}
```
`srcAddress` is the address the stream is pulled from. This is usually a RTSP source, but also RTMP and HLS sources are allowed. `optionalOutputAddress` is the RTMP address where the stream should be repeated to. Currently only RTMP targets are allowed.

```json
"options": {
    "rtspTcp": true
}
```
`rtspTcp` is whether to enable RTSP transport via TCP in contrast to UDP.

```json
"states": {
    "repeatToLocalNginx": {
        "type": "connected"
    },
    "repeatToOptionalOutput": {
        "type": "connected"
    }
}
```
`repeatToLocalNginx` defines the state of the stream from `srcAddress`. It can have the values `connected` and `disconnected`. `connected` means, the Restreamer is pulling the stream. The same for `repeatToOptionalOutput`, the RTMP target the stream is repeated to.

```json
"userActions": {
    "repeatToLocalNginx": "start",
    "repeatToOptionalOutput": "start"
}
```
The `userActions` are the desired state of the source and target stream. Set it to `start` and the type to `connected` in the `states` section to autostart pulling or pushing the stream. If you do not want to autostart, then set the user action to `stop` and the state to `disconnected`.

Please keep in mind that this format of the file may change in future releases.
