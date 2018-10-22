---
title: Raspicam
---

###### [User Guides](../docs/guides-index.html) > Raspicam 
## Raspicam

---
**Datarhei Hint ☺ ►**  for Hypriot-Users   
The Hypriot-Image has no raspi-config to activate the camera. You have to install this first:

```sh
# cd /tmp
# apt-get update && apt-get install -y alsa-utils
# wget http://archive.raspberrypi.org/debian/pool/main/r/raspi-config/raspi-config_20160108_all.deb
# dpkg -i raspi-config_20160108_all.deb
```

*Source <a target= "_blank" href="https://github.com/snubbegbg">here</a>*

---

1. Setup your RaspiCam  
  `# raspi-config` -> "Enable Camera" -> "Finished" and then reboot  
2. Show the running Docker image by name and ID   
  `# docker ps`
3. Kill the running container:  
  `# docker kill ContainerName && docker rm ContainerName`
4. Start the Docker image  
   * Raspberry Pi 1:   
  `# docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -e "MODE=RASPICAM" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db -v /opt/vc:/opt/vc --privileged datarhei/restreamer-armv6l:latest`
   * Raspberry Pi 2:   
  `# docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -e "MODE=RASPICAM" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db -v /opt/vc:/opt/vc --privileged datarhei/restreamer-armv7l:latest`
5. Put the address into the "RTMP/RTSP Video Source" field  
   `rtmp://127.0.0.1/live/raspicam.stream`   
   
---