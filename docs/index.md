---
title: Documentation
---

Welcome to the Restreamer documentation, which helps you to install and configure Restreamer for your device.
All other questions about the streaming microcosm and technical facts are explained [in our Wiki](../wiki). 

Restreamer is built on [Docker](../wiki/docker.html). If you don't know Docker, don't be scared. Extensive knowledge isn't required!
{: .notice--info}

## Requirements

* An IP camera, USB camera, or RaspiCam
* A x86_64 or ARM computer running Windows, macOS, or Linux
* A working Docker installation

## Quick start

Assuming that you already have a device with Docker running on it, fire up the console and type:

```sh
$ docker run -it --rm -p 8080:8080 datarhei/restreamer:latest
```

Choose the image `datarhei/restreamer-armv7l:latest` if you are running it on a device with an ARM CPU in 32bit mode.
{: .notice--warning}

Open your browser, load the URL `http://your-device-ip:8080/` and enter `admin` as username and `datarhei` as password. Now you
are ready to stream!

## Streaming Guides

Check our user guides on the left to find out more about how to stream your [IP Camera](guides-ipcam-rtsp.html),
your [Raspicam](guides-raspicam.html), or your [USB Camera](guides-usb-camera.html).

## Installation Guides

For more details on how to run Restreamer on your device, please refer to our installation guides:

- [Windows](installation-osx-windows.html)
- [macOS](installation-osx-windows.html)
- [Linux](installation-linux-64.html)
- [Raspberry Pi](installation-linux-arm.html)
- [Odroid](installation-linux-arm.html)

## Configuration Guides

You can customize your Restreamer by setting different [environment variables](references-environment-vars.html), e.g. for changing the
[login for the GUI](references-environment-vars.html#rs_username), the [snapshot interval](references-environment-vars.html#rs_snapshot_interval),
or tampering [the audio track](references-environment-vars.html#rs_audio) of your stream.
