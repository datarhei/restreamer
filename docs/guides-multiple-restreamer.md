---
title: Multiple Restreamer
---

It is possible to run multiple Restreamer on one device, i.e. in order to restream more than one camera or stream on one device.
This guide will show you what to take care of such that the different Restreamer containers don't step on each others feet.

There three things to consider when starting the Docker images:

* each containter needs a different name
* each containter needs a different port
* each containter needs a different volume

## Example

Imagine that you want to run 4 Restreamer instances on one device. The Docker command lines could
look like this

```sh
docker run -d --restart always -e "RS_USERNAME=..." -e "RS_PASSWORD=..." --name camera1 -p 8091:8080 -v /camera1/db:/mnt/restreamer/db datarhei/restreamer-armv7l:latest
docker run -d --restart always -e "RS_USERNAME=..." -e "RS_PASSWORD=..." --name camera2 -p 8092:8080 -v /camera2/db:/mnt/restreamer/db datarhei/restreamer-armv7l:latest
docker run -d --restart always -e "RS_USERNAME=..." -e "RS_PASSWORD=..." --name camera3 -p 8093:8080 -v /camera3/db:/mnt/restreamer/db datarhei/restreamer-armv7l:latest
docker run -d --restart always -e "RS_USERNAME=..." -e "RS_PASSWORD=..." --name camera4 -p 8094:8080 -v /camera4/db:/mnt/restreamer/db datarhei/restreamer-armv7l:latest
```

Note that each container has a different name
```sh
... --name camera1 ...
... --name camera2 ...
... --name camera3 ...
... --name camera4 ...
```

To give an explicit name is not required, but makes it easier to work with the containers.

Each container has to map to a different port on your device
```sh
... -p 8091:8080 ...
... -p 8092:8080 ...
... -p 8093:8080 ...
... -p 8094:8080 ...
```

The instance called `camera1` will be reachable on `http://localhost:8091/`, the instance `camera2` on `http://localhost:8092/`, and so on. If the ports
are not different, also Docker would complain.

Most importantly is that you map different directories to the `/mnt/restreamer/db` volume
```sh
... -v /camera1/db:/mnt/restreamer/db ...
... -v /camera2/db:/mnt/restreamer/db ...
... -v /camera3/db:/mnt/restreamer/db ...
... -v /camera4/db:/mnt/restreamer/db ...
```
If you map always the same directory, the different Restreamer instances would overwrite the configuration files of each other. The configuration file
contains all required information about the currently connected streams. In case of a restart the Restreamer can automatically reconnect to the previously
connected streams.

With different directories mounted to the volume, each Restreamer has its own configuration file that will not get overwritten by the other Restreamers.
