---
title: SD Card
---

When you use an embedded device or any other device with Restreamer running from a SD card, then you should be aware
that Restreamer is causing quite some I/O while streaming and logging (depending on the log level you set).

This high I/O might wear out your SD card very quickly which may lead to errors or crashes during operation. Eventually
you have to replace the SD card and restore your whole setup.

Here are some countermeasures to avoid unecessary wear-out of your SD card.

## tmpfs

Restreamer uses the directory `/tmp/hls` to store the HLS chunks and generates a lot I/O operations during streaming. Mount this directory
as a `tmpfs` (RAM drive) by adding this to you Docker command line: `docker run ... --tmpfs /tmp/hls ...`. Docker will then create a
`tmpfs` and mount it to `/tmp/hls`.

## Logging

Extensive logging may also cause some significant I/O on the SD card. You can either silence Restreamer by setting the [environment
variable](references-environment-vars.html) `RS_LOGLEVEL` to `0` or tell Docker not ot log anything: `docker run ... --log-driver none ...`.

If you want to keep logs from Docker and Restreamer, then you can tell Docker to write the logs to syslog:
`docker run ... --log-driver syslog ...`. Make sure that the directory where you store the syslogs is mounted to a `tmpfs` on your
host system.

Find out more about the different [Docker log drivers](https://docs.docker.com/config/containers/logging/configure/).