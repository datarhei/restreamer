---
title: Encoding
---

As of Restreamer version 0.4.0 you can choose to encode the video stream that you are pulling from the camera to [H.264](../wiki/h264.html). 
Encoding of the video streams enables you to use stream sources that are not encoded in H.264 or to re-encode a video stream that is
already encoded in H.264, but you want to change e.g. the bitrate.

The encoding is always towards H.264 as this is the most compatible codec for HLS streaming and the only useful choice for pushing
the stream to a RTMP server, e.g. Youtube, Twitch, ...

## Default Settings

By default, the settings for the video codec is `copy` and for the audio codec `auto`. These are also the default
settings for the Restreamer versions before 0.4.0.

If you set the `RS_AUDIO` environment variable on the command line, its value will be translated to the following settings:

Value | Codec | Preset | Bitrate | Channels | Sampling
------|-------|--------|---------|----------|---------
`auto` | `auto` | - | - | - | -
`none` | `none` | - | - | - | -
`silence` | `aac` | `silence` | `64 Kbit/s` | `mono` | `44100`
`aac` | `aac` | `encode` | `64 Kbit/s` | `copy`? | `copy`?
`mp3` | `mp3` | `encode` | `64 Kbit/s` | `copy`? | `copy`?

## Video Setting

The default video codec setting is `copy`. This means that the video stream from the source is used unmodified for the player and for
optionally pushing it to an external RTMP server.

### H.264

Choose `H264` as video codec to enforce an encoding of the video stream from the source. This will reveal more options

Encoding the video stream may require significant CPU resources.
{: .notice--warning}