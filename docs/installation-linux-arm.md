---
title: Installation Raspberry Pi / ARM
---

## Requirements

* Raspberry Pi 1 / Pi 2 / Pi 3 or Odroid U3
* ARM CPU that supports the `arm32v6`, `arm32v7`, or `arm64v8` architecture

## Raspberry Pi 1 / Pi 2 / Pi 3

Verified with:

* [Raspberry Pi 1 Model B+](https://www.raspberrypi.org/products/model-b-plus/)
* [Raspberry Pi 2 Model B](https://www.raspberrypi.org/products/raspberry-pi-2-model-b/)
* [Raspberry Pi 3 Model B+](https://www.raspberrypi.org/products/raspberry-pi-3-model-b-plus/)

Instructions:

1. Install the latest [Raspian Stretch](https://www.raspberrypi.org/downloads/raspbian/) image on the SD card

2. Login to the Pi and install Docker CE according to the [Debian install instructions](https://docs.docker.com/install/linux/docker-ce/debian/#install-from-a-package)

## Odroid U3

Verified with:

* [Odroid U3](http://www.hardkernel.com/main/products/prdt_info.php?g_code=g138745696275) (discontinued)

Instructions:

1. Install the latest [Ubuntu](https://com.odroid.com/sigong/nf_file_board/nfile_board.php) image on the SD card

2. Login to the Odroid and install Docker CE according to the [Unbuntu install instructions](https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-from-a-package)


## Common Steps

3. Start the Restreamer   
   ```sh
   $ docker run -d --restart always \
        --name restreamer \
        -e "RS_USERNAME=admin" -e "RS_PASSWORD=datarhei" \
        -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db \
        --tmpfs /tmp/hls \
        datarhei/restreamer-armv7l:latest
   ```
4. Browse to http://your-device-ip:8080

The default login is:

* Username: `admin`
* Password: `datarhei`

It is highly recommended to change the username and password.
{: .notice--warning}

## Important customizations

It is recommended to change the username and password. In order to change them you have to set the respective [environment variables](references-environment-vars.html)
in the docker command:

```sh
$ docker run ... -e "RS_USERNAME=YOUR_USERNAME" -e "RS_PASSWORD=YOUR_PASSWORD" ...
```

## Description of the command

#### `-d`

Detach the container. This means that the container will run in the background. You can run it interactively in the foreground with
`-it` instead of `-d`. To stop the Restreamer in detached mode, type `docker stop restreamer`. In interactive mode just hit `Ctrl-C` to
stop the Restreamer.

#### `--name restreamer`

Gives the container the name `restreamer`. This name can be used in other docker commands to control the container. In order to
stop the Restreamer, type `docker stop restreamer`. While the Restreamer is running you can log in into the container with `docker exec -it restreamer /bin/bash`.

#### `--restart always`

In case the Restreamer crashes, Docker will automatically restart the Restreamer.

#### `-e "RS_USERNAME=..." -e "RS_PASSWORD=..."`

Set values for the environment variables `RS_USERNAME` and `RS_PASSWORD`. See a description of all known [environment variables](references-environment-vars.html).

#### `-p 8080:8080`

Bind the port 8080 of the device to the port 8080 of the Restreamer. With this you can connect with your browser to the Restreamer GUI.
If you want to us another port, change it to e.g. `-p 31000:8080`.

#### `--tmpfs /tmp/hls`

This will mount the directory `/tmp/hls` as a disk in the RAM of the container. This directory holds the chunks for the HLS chunks and has a lot I/O
operations on it. By mounting it to a `tmpfs` you avoid [wearing out the SD card](guides-sdcard.html) in your device.

#### `-v /mnt/restreamer/db:/restreamer/db`

The Restreamer stores the current state in the directory `/restreamer/db` inside the container. This command maps the directory `/mnt/restreamer/db`
of your device into the container. With this the state can be preserved in case the Restreamer needs to be restarted. If you want to store
the state in a different directory on your device, change it to e.g. `-v /tmp/restreamer:/restreamer/db`

#### `datarhei/restreamer-armv7l:latest`

This is the docker image of the lastest Restreamer on the Docker Hub. Docker will check if the image is locally available
and download it if it is not available or a newer image is available.
