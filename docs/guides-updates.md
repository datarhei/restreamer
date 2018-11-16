---
title: Update Restreamer
---

To update the Restreamer without loosing the data of the input fields it is required to start the Docker-Image with `-v /path/to/local:/restreamer/db` ([details](installation-linux-64.html#description-of-the-command)) or a activated volume in Kitematic (Docker Toolbox) ([details](installation-osx-windows.html#important-customizations))!

* [OSX / Windows](#osx-windows)
* [Linux](#linux)


## OSX / Windows

1. Open Kitematic and click in the running Restreamer on stop:
![Kitematic](../img/references-updates-stop.png)
2. Delete the old image:
![Kitematic](../img/references-updates-delete.png)
3. Start the image again: 
![Kitematic](../img/references-updates-create.png)
4. Enable the volume again:
![Kitematic](../img/references-updates-db-restore.png)
5. Insert your enviroments again ([more here](installation-osx-windows.html#important-customizations))

Done! 

## Linux

1.
  Stop and remove the running Restreamer:   
  ```sh
  docker stop restreamer && docker rm restreamer
  ```
2.
  Download the new image. Please do not forget to add the right arm-tag if used:
  ```sh
  # for amd64 CPUs
  docker pull datarhei/restreamer:latest  

  # for arm32v7 CPUs Raspberry Pi 2, Pi 3, ...
  docker pull datarhei/restreamer-armv7l:latest

  # for arm64v8 CPUs
  docker pull datarhei/restreamer-aarch64:latest

  # for arm32v6 CPUs
  docker pull datarhei/restreamer-armv6l:latest
  ```
3.
  Start the Restreamer again (please remember to modify if you have your own configuration):   
  ```sh
  docker run -d --restart always \
    --name restreamer \
    -e "RS_USERNAME=YOUR-USERNAME" -e "RS_PASSWORD=YOUR-PASSWORD" \
    -p 8080:8080 \
    -v /mnt/restreamer/db:/restreamer/db \
    datarhei/restreamer:latest
  ```

Done!
