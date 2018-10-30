---
title: Raspicam
---


## Enable RaspiCam

Execute the `raspi-config` command and then select `Enable Camera`. Leave `raspi-config` and reboot your Raspberry Pi.

## Configure Restreamer

For Restreamer to be able to access the RaspiCam, you have to stop the Restreamer if it is currently running:

```sh
docker stop restreamer
docker rm restreamer
```

Restart the Restreamer with the RaspiCam mode enabled:

```sh 
docker run -d --restart always \
    --name restreamer \
    -e "RS_USERNAME=..." -e "RS_PASSWORD=..." -e "MODE=RASPICAM" \
    -p 8080:8080 \
    -v /mnt/restreamer/db:/restreamer/db \
    -v /opt/vc:/opt/vc \
    --privileged \
    datarhei/restreamer-armhf:latest
```

In order to stream what your RaspiCam is recording, you have to put `rtmp://127.0.0.1/live/raspicam.stream` into
the "RTMP/RTSP/HLS Video Source" field and press "Start".

<img src="../img/guides-raspicam-url.png" width="95%">

## raspi-config is missing

If you are running a Linux on your Raspberry Pi that is missing the `raspi-config` program, you have to install it first:

```sh
cd /tmp
sudo apt-get update && apt-get install -y alsa-utils
wget https://archive.raspberrypi.org/debian/pool/main/r/raspi-config/raspi-config_20180406+1_all.deb
sudo dpkg -i raspi-config_20180406+1_all.deb
```

Please check `https://archive.raspberrypi.org/debian/pool/main/r/raspi-config/` for the lastest available version of `raspi-config`.

[Source](https://github.com/snubbegbg/install_raspi-config)