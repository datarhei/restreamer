---
title: USB-Cameras
---

## USB-Cameras

1. Connect you USB-Camera
2. Run Docker-Image:   
   `docker run  -d --name restreamer -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db --device /dev/video0 datarhei/restreamer-armv7l:latest`
3. Install additional software:   
   `docker exec -it restreamer /bin/bash`   
   * apt-get install v4l-utils libv4l-0
   * ffmpeg -f v4l2 -r 25 -s 1280x720 -i /dev/video0 -f flv rtmp://127.0.0.1:1935/live/usb.stream
4. Now you can put the adress   
   `rtmp://127.0.0.1/live/usb.stream`
   into the "RTMP/RTSP Video Source" field

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.