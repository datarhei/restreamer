---
title: guides-setup-video-device 
---
# Setup your video device
For the integration of an H.264 -enabled network camera you need the address, on which you can retrieve the video live stream from the camera.
Take a look at your camera manual and look after RTSP/RTP.  

For this purpose, please look in the manual of your camera and search for " RTSP " . Alternatively you can find already many templates on the Soleratec company website: [https://www.soleratec.com/support/rtsp/rtsp_listing](https://www.soleratec.com/support/rtsp/rtsp_listing)  

**Two examples:**  
AXIS: rtsp://ip-address:554/axis-media/media.amp  
Samsung: rtsp://ip-address:554/profile2/media.smp  

*You can finde the IP address of your camera, if necessary, in the DHCP table of your router. You can  use a LAN scanner like http://angryip.org/ which displays you all the devices on your network too if you are unsure what the correct IP is.*  

Next step: Put your rtsp-link in “RTMP/RTSP Video Source” field in the Datarhei/Restreamer-Web-GUI.  

![UI-Preview](../img/setup-video-device.png)  

Last step: Start process  

Once the process has been successfully established , you can open the Datarhei/Restramer player and, if necessary, forward the stream to an external provider the player.  

---

Want to talk to us? Write email open@datarhei.org, go to [Support](../support.html) or choose a nickname and join us on <a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">#datarhei webchat on freenode</a>.

If you're having a weird problem while developing, see [Known Issues](https://github.com/datarhei/small-restreamer-internal/issues/). 
