---
title: RaspiCam
---

## RaspiCam

1. Setup your RaspiCam
2. Start the Docker-Image: `docker run -d -v /mnt/restreamer/db:/restreamer/db --restart always --name restreamer -p 8080:8080 -v /opt/vc:/opt/vc --privileged -e "MODE=RASPICAM" datarhei/restreamer-armv7l:latest`
3. Put the adress `rtmp://127.0.0.1/live/raspicam.stream` into the "RTMP/RTSP Video Source" field

---

Want to talk to us? Write email open@datarhei.org, go to [Support](../support.html) or choose a nickname and join us on <a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">#datarhei webchat on freenode</a>.

If you're having a weird problem while developing, see [Known Issues](https://github.com/datarhei/small-restreamer-internal/issues/). 
