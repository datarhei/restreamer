---
title: RaspiCam
---

## RaspiCam

If you use the Hypriot-Image for Docker please read the hint under the docker-guide.

1. Setup your RaspiCam   
  `# raspi-config` -> "Enable Camera" -> "Finised" and then reboot
2. Start the Docker-Image
   * Raspberry Pi 1:   
  `# docker run -d -v /mnt/restreamer/db:/restreamer/db --restart always --name restreamer -p 8080:8080 -v /opt/vc:/opt/vc --privileged -e "MODE=RASPICAM" datarhei/restreamer-armv6l:latest`
   * Raspberry Pi 2:   
  `# docker run -d -v /mnt/restreamer/db:/restreamer/db --restart always --name restreamer -p 8080:8080 -v /opt/vc:/opt/vc --privileged -e "MODE=RASPICAM" datarhei/restreamer-armv7l:latest`
3. Put the address   
   `rtmp://127.0.0.1/live/raspicam.stream`   
   into the "RTMP/RTSP Video Source" field

**Datarhei Hint ☺ ►**  for Hypriot-Users   
The Hypriot-Image has no raspi-config to activate the camera. You have to install this first:

```sh
# cd /tmp
# wget http://archive.raspberrypi.org/debian/pool/main/r/raspi-config/raspi-config_20150131-5_all.deb
# apt-get update && apt-get install -y libnewt0.52 whiptail parted triggerhappy lua5.1
# dpkg -i raspi-config_20150131-5_all.deb
```
*Source [here](https://github.com/snubbegbg)*

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.
