---
title: Architecture
---
###### [Development](../docs/development-index.html) > Architecture
# Architecture

The Datarhei/Restreamer basically consists of four different applications:

* Frontend with Angular and Node.js for process management
* NGINX web server with the RTMP-Module
* FFmpeg as swiss army knife for video processing
* Clappr video player

---
![UI-Preview](../img/architecture.png)

---
### Processes:

1. The application provides the user interface / HTTP API and starts the NGINX web server using the supplied config  (/restreamer/config/nginx.conf)  
2. FFmpeg fetches the camera stream and forwards it to the local running NGINX RTMP server rtmp://127.0.0.1:1935/hls/live.stream  
3. NGINX-RTMP delivers the camera stream to the address: http://...:8080/hls/live.stream.m3u8
4. Clappr videoplayer calls the camera stream via HTTP (HLS)
5. Additonally, FFmpeg takes the local NGINX-RTMP video stream rtmp://127.0.0.1:1935/hls/live.stream and pushes the source to 
the registered favourite address

---