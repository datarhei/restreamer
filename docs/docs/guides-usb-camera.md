---
title: USB-Cameras
---

## USB-Cameras

1. Connect you USB-Camera
2. Run Docker-Image: `docker run  -d --name restreamer -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db --device /dev/video0 datarhei/restreamer-armv7l:latest`
3. Install additional software: `docker exec -it restreamer /bin/bash`   
   * apt-get install v4l-utils libv4l-0
   * ffmpeg -f v4l2 -r 25 -s 1280x720 -i /dev/video0 -f flv rtmp://127.0.0.1:1935/live/usb.stream
4. Now you can put the adress `rtmp://127.0.0.1/live/usb.stream` into the "RTMP/RTSP Video Source" field

---

Want to talk to us? Write email open@datarhei.org, go to [Support](../support.html) or choose a nickname and join us on <a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">#datarhei webchat on freenode</a>.

If you're having a weird problem while developing, see [Known Issues](https://github.com/datarhei/small-restreamer-internal/issues/). 
