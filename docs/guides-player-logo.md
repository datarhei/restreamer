---
title: Logo
---
###### [User Guides](../docs/guides-index.html) > Player Logo

### How to add/change the watermark in your video player

Datarhei/Restreamer uses the open source video player clappr.io full documentation on clappr player and more things
 <a target= "_blank" href="https://github.com/clappr/clappr/#watermark">here</a>.

This guide helps you to change or add a watermark as your company logo to the Datarhei/Restreamer video stream.
**Watermark is just working with the Datarhei/Restreamer clappr video player delivered through the direct streaming iframe!**

Requirement for this guide is a full installed and working Datarhei/Restreamer connected to the internet and with power on!

1. SSH  login to your Datarhei/Restreamer (open CLI in Kitematic (Docker Toolbox))

2. Switch in the Docker container to add the Player values with:

```sh
# docker exec -it restreamer /bin/bash
```


3. You are now in the instance of the Docker Imaage. It is like a own shell. Change directory with the command:  

```sh 
# cd /restreamer/bin/webserver/public
```

4. Locate and open index.html with editor of your choice (apt-get ist working in this Docker Image install your favourite editor if it is not already installed)

5. Locate in opened index.html:  

```sh
var player = new Clappr.Player({  
source: "http://" + window.location.hostname + ":" + window.location.port + "/hls/live.stream.m3u8",  
parentId: "#player",  
baseUrl: '/libs/clappr/dist/',  
poster: "images/live.jpg",  
mediacontrol: {seekbar: "#3daa48", buttons: "#3daa48"},  
height: "100%",  
width: "100%"  
# }); 
```  


6. Add this two lines and change the link to your imagesource and the watermark link. If you dont need a link dont add this line. Linking will be automaticly disabled:   

```sh 
# watermark: "http://linktoyourimage.png", position: 'top-right',   
#`watermarkLink: "http://datarhei.org/"
```

7. should look now:  

```sh 
# var player = new Clappr.Player({  
#  source: "http://" + window.location.hostname + ":" + window.location.port + "/hls/live.stream.m3u8",  
#  parentId: "#player",  
#  baseUrl: '/libs/clappr/dist/',  
#  poster: "images/live.jpg",  
#  mediacontrol: {seekbar: "#3daa48", buttons: "#3daa48"},  
#  watermark: "http://linktoyourimage.png", position: 'top-right',  
#  watermarkLink: "http://datarhei.org/",  
#  height: "100%",  
#  width: "100%"  
# }); 
```


8. Watermark positions can be bottom-left, bottom-right, top-left and top-right. It is possible to add 4 watermarks in one player two at the upper and lower corners



---

---