---
title: USB-Cameras
---
###### [User Guides](../docs/guides-index.html) > USB-Camera
## USB-Camera

---
**Datarhei Hint ☺ ►** for ARMv6l / ARMv7l (e.g. Raspberry Pi)
It works with USB but it is not funny on ARM CPU cause of very low FPS rates. On more powerfull systems it is working much better with higher FPS results.  

---
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

---
**Datarhei Hint ☺ ►**   
The first Restreamer start will take a while because he has to install some packages for the usb-driver

---