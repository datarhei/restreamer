---
title: USB camera
---

On ARMv6l / ARMv7l (e.g. Raspberry Pi), USB camera work but only with a very low framerate. It will work with a higher
framerate on more powerfull systems. Because the video stream from an USB camera is usually not encoded in H.264 (which is a
requirement for Restreamer), it needs to be encoded. This encoding requires significant CPU resources.

Your USB camera needs to be v4l compliant in order to work with Restreamer.
{: .notice--info}

## Configure Restreamer

Connect the USB camera. For Restreamer to be able to access the USB camera, you have to stop the Restreamer if it is currently running:

```sh
docker stop restreamer
docker rm restreamer
```

Restart the Restreamer with the USB camera mode enabled:

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

![USB camera](../img/guides-usbcam-url.png)

The first Restreamer start will take a while because it has to install some packages for the USB driver.
{: .notice--info}

## Environment Variables

These environment variables enable you to control the Raspberry Pi camera.

| Name | Default | Description |
|------|---------|-------------|
| `RS_USBCAM_DEVICE` | `/dev/video0` | The video device that you also pass in the command line for Docker. |
| `RS_USBCAM_FPS` | `25` | Set the framerate. |
| `RS_USBCAM_GOP` | `50` | Set the GOP size. Usually the framerate multiplied by 2 is a good value. |
| `RS_USBCAM_BITRATE` | `5242880` | Bitrate of the video stream in bit/s. |
| `RS_USBCAM_WIDTH` | `1920` | Video stream width in pixels. |
| `RS_USBCAM_HEIGHT` | `1080` | Video stream height in pixels. |

Change the defaults of these environment variable with care and make sure that you know what you are doing.