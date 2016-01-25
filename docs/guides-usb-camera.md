---
title: USB-Cameras
---

## USB-Cameras

1. Connect you USB-Camera
2. Show the running Docker image by name and ID   
  `# docker ps`
3. Kill the running container:  
  `# docker kill ContainerName && docker rm ContainerName`
4. Start the Docker image  
   * Raspberry Pi 1:   
  `# docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -e "MODE=USBCAM" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db --device /dev/video0 datarhei/restreamer-armv6l:latest`
   * Raspberry Pi 2:   
  `# docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -e "MODE=USBCAM" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db --device /dev/video0 datarhei/restreamer-armv7l:latest`
5. Put the address into the "RTMP/RTSP Video Source" field  
   `rtmp://127.0.0.1/live/usbcam.stream` 

**Datarhei Hint ☺ ►**   
The first Restreamer start will take a while because he has to install some packages for the usb-driver

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.
