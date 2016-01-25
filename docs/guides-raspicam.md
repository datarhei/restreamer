---
title: RaspiCam
---

## RaspiCam

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
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.
