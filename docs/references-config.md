---
title: Configuration
---

This is the basic configuration file for Restreamer and it is is located in `conf/live.json`. It hosts settings for different
internal aspects of the Restreamer.

Usually you don't need to change these values. You can safely control them with the available [enviroment variables](references-environment-vars.html).
{: .notice--success}

Please don't change these settings if you don't understand the internals of Restreamer.
{: .notice--danger}

## Name

```json
{
    "name": "live"
}
```
If you want to add more than one config, you have to change this name to "dev" or something you like and start the application with the env. "MODE=dev" to select the other one.

## JSON database

```json
{
    "jsondb": "db/v1"
}
```

The path to the JSON database. More about the [JSON database](https://www.npmjs.com/package/node-jsondb).

## Auth

```json
{
    "auth": {
        "username": "admin",   
        "password": "datarhei",
        "token": ""
    }
}
```
Default values for the login credentials if they are not set by the environment variables.
More about the [environment variables](references-environment-vars.html).

## FFmpeg options

```json
{
    "ffmpeg": {
        "options": {
            "native_h264_native_audio": {
                "outputOptions": [
                    "-codec copy",
                    "-map 0:v",
                    "-map 0:a",
                    "-map_metadata -1",
                    "-metadata application=datarhei/Restreamer",
                    "-metadata server=NGINX-RTMP",
                    "-f flv"
                ]
            },
            "native_h264_no_audio": {
                "outputOptions": [
                    "-vcodec copy",
                    "-an",
                    "-map 0:v",
                    "-map_metadata -1",
                    "-metadata application=datarhei/Restreamer",
                    "-metadata server=NGINX-RTMP",
                    "-f flv"
                ]
            },
            "native_h264_native_aac": {
                "outputOptions": [
                    "-codec copy",
                    "-map 0:v",
                    "-map 0:a",
                    "-bsf:a aac_adtstoasc",
                    "-map_metadata -1",
                    "-metadata application=datarhei/Restreamer",
                    "-metadata server=NGINX-RTMP",
                    "-f flv"
                ]
            },
            "native_h264_silence_aac": {
                "input": "anullsrc=r=44100:cl=mono",
                "inputOptions": [
                    "-f lavfi"
                ],
                "outputOptions": [
                    "-vcodec copy",
                    "-acodec aac",
                    "-b:a 0k",
                    "-map 0:v",
                    "-map 1:a",
                    "-shortest",
                    "-map_metadata -1",
                    "-metadata application=datarhei/Restreamer",
                    "-metadata server=NGINX-RTMP",
                    "-f flv"
                ]
            },
            "native_h264_transcode_aac": {
                "outputOptions": [
                    "-vcodec copy",
                    "-acodec aac",
                    "-b:a 64k",
                    "-map 0:v",
                    "-map 0:a",
                    "-map_metadata -1",
                    "-metadata application=datarhei/Restreamer",
                    "-metadata server=NGINX-RTMP",
                    "-f flv"
                ]
            },
            "native_h264_transcode_mp3": {
                "outputOptions": [
                    "-vcodec copy",
                    "-acodec libmp3lame",
                    "-b:a 64k",
                    "-map 0:v",
                    "-map 0:a",
                    "-map_metadata -1",
                    "-metadata application=datarhei/Restreamer",
                    "-metadata server=NGINX-RTMP",
                    "-f flv"
                ]
            },
            "global": {
                "inputOptions": [
                    "-stats",
                    "-loglevel quiet",
                    "-err_detect ignore_err"
                ]
            },
            "rtsp-tcp": {
                "inputOptions": [
                    "-rtsp_transport tcp"
                ]
            },
            "snapshot": {
                "outputOptions": [
                    "-vframes 1"
                ]
            }
        }
    }
}
```

Preset for FFmpeg. The selected presets depend on the incoming stream and the environment variable for the audio setting. By default, it is assumed, that
the stream has no audio (`native_h264_no_audio`). If there is an audio track, `native_h264_native_audio` will be selected. In case the audio track
is encoded with AAC, then the preset `native_h264_native_aac` gets selected (to fix possible errors in the audio stream). If the audio stream is
not supported by the FLV container, it will be transcoded to AAC with the preset `native_h264_transcode_aac`. 

The selection of the incoming video stream preset can be influenced by the `RS_AUDIO` [enviroment variable](docs/references-environment-vars.html).
{: .notice--info}
        
## FFmpeg monitor

```json
{
    "ffmpeg": {
        "monitor": {
            "restart_wait": "6000",  
            "stale_wait": "60000"
        }
    }
}
```

The running FFmpeg process is monitored. If the process stops, it will be automatically restarted after `restart_wait` milliseconds. If the process
is stale, i.e. no incoming data is processed, FFmpeg will be automatically restart after `stale_wait` milliseconds.

## NGINX

```json
{
    "nginx": {
        "command": "/usr/local/nginx/sbin/nginx",
        "args": [
            "-c",
            "/restreamer/conf/nginx.conf"
        ],
        "streaming": {
            "ip": "127.0.0.1",
            "rtmp_port": "1935",
            "rtmp_hls_path": "/hls/",
            "http_port": "8080",
            "http_health_path": "/ping"
        }
    },
}
```

Preset for the NGINX process. The values in the `streaming` section reflect what is configured in the `nginx.conf` file.

Don't change any of these parameter without changing them accordingly in the `nginx.conf` file.
{: .notice--danger}

## Environment variables
```json
{
    "envVars": [
        {
            "name": "RS_NODEJS_PORT",
            "alias": [
                "NODEJS_PORT"
            ],
            "type": "int",
            "defaultValue": "3000",
            "required": false,
            "description": "Webserver port of application."
        },
        {
            ...
        }
    ]
}
```

Definition of the available environment variables, their type, default value, and description.
