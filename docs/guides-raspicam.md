---
title: RaspiCam
---

## RaspiCam

1. Setup your RaspiCam
2. Start the Docker-Image:   
  `# docker run -d -v /mnt/restreamer/db:/restreamer/db --restart always --name restreamer -p 8080:8080 -v /opt/vc:/opt/vc --privileged -e "MODE=RASPICAM" datarhei/restreamer-armv7l:latest`
3. Put the adress   
   `rtmp://127.0.0.1/live/raspicam.stream`   
   into the "RTMP/RTSP Video Source" field

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.