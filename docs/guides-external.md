---
title: External Streaming Server
---

With Restreamer you can forward your stream to an external video streaming server that support either [RTMP](../wiki/rtmp.html) or [HLS](../wiki/hls-http.html).

## RTMP

RTMP is most commonly used by streaming providers, such as YouTube, Facebook, and so on. Please follow these guides:

- [General RTMP streaming](guides-external-rtmp.html)
- [YouTube](guides-youtube.html)
- [Twitch](guides-twitch.html)

## HLS

HLS is not so common for pushing the stream to an external streaming servers, but it is practical if you don't want to setup
an RTMP server because the stream is pushed via normal HTTP POST to the server.

Please follow the guide for [external HLS streaming](guides-external-hls.html).
