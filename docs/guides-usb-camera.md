---
title: USB-Camera
---

On ARMv6l / ARMv7l (e.g. Raspberry Pi), it works with USB but only with a very low framerate. It will work with a higher framerate on more powerfull systems.
{: .notice--warning}

## Configure Restreamer

Connect the USB-Camera. For Restreamer to be able to access the USB-Camera, you have to stop the Restreamer if it is currently running:

```sh
docker stop restreamer
docker rm restreamer
```

Restart the Restreamer with the USB-Camera mode enabled:

```sh 
docker run -d --restart always \
    --name restreamer \
    -e "RS_USERNAME=..." -e "RS_PASSWORD=..." -e "MODE=USBCAM" \
    -p 8080:8080 \
    -v /mnt/restreamer/db:/restreamer/db \
    --device /dev/video0 \
    datarhei/restreamer-armhf:latest
```

In order to stream what your USB-Camera is recording, you have to put `rtmp://127.0.0.1/live/usbcam.stream` into
the "RTMP/RTSP/HLS Video Source" field and press "Start".

<img src="../img/guides-usbcam-url.png" width="95%">

The first Restreamer start will take a while because it has to install some packages for the USB driver.
{: .notice--info}