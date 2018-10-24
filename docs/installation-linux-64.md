---
title: Installation Linux
---

## Requirements

* 64bit Intel or AMD CPU
* A supported Linux distribution

Currently officially supported Linux distributions are:

* [CentOS](https://docs.docker.com/install/linux/docker-ce/centos/)
* [Debian](https://docs.docker.com/install/linux/docker-ce/debian/)
* [Fedora](https://docs.docker.com/install/linux/docker-ce/fedora/)
* [Ubuntu](https://docs.docker.com/install/linux/docker-ce/ubuntu/)

Other Linux distributions that are derived from the ones listed above may also be supported. Please check the respective installation guides.
{: .notice--info}


## Installation

1. Download and install Docker CE. Please refer to the guides that are listed above.

2. Start the Restreamer  
   ```shell
   # docker run -d --name restreamer --restart always -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer:latest
   ```

3. Browse to http://your-device-ip:8080

The default login is:

* Username: `admin`
* Password: `datarhei`

It is highly recommended to change the username and password.
{: .notice--warning}

## Important customizations

It is recommended to change the username and password. In order to change them you have to set the respective [environment variables](references-environment-vars.html)
in the docker command:

```shell
# docker run ... -e "RS_USERNAME=YOUR_USERNAME" -e "RS_PASSWORD=YOUR_PASSWORD" ...
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

#### `-v /mnt/restreamer/db:/restreamer/db`

The Restreamer stores the current state in the directory `/restreamer/db` inside the container. This command maps the directory `/mnt/restreamer/db`
of your device into the container. With this the state can be preserved in case the Restreamer needs to be restarted. If you want to store
the state in a different directory on your device, change it to e.g. `-v /tmp/restreamer:/restreamer/db`

#### `datarhei/restreamer:latest`

This is the docker image of the lastest Restreamer on the Docker Hub. Docker will check if the image is locally available
and download it if it is not available or a newer image is available.
