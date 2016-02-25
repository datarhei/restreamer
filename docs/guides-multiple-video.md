---
title: guides-multiple-video
---
# How to add multiple video sources to Datarhei/Restreamer

If you have more than one video source you have to start a second Docker image. 
In the example we start the second Docker image with new parameters to run it simultanious with the default Docker Image:  
* change the default name restreamer to restreamer2   
* change the default port 8080 to 8090  
- add a second db named db2  

## Content

* [Raspberry Pi1 (ARMv6l)](#raspberry-pi1-(armv6l))    
* [Raspberry Pi2 (ARMv7l)](#armv7l)    
* [Linux 64 bit](#linux-64-bit)    
* [MacOS and Windows](#macos-and-windows)    

---

## Raspberry Pi1 (ARMv6l)  

1. SSH Login to your ARMv6l device
1. Start first image "restreamer" if not already running:  
```
docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=admin" -e "RESTREAMER_PASSWORD=datarhei" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db --privileged datarhei/restreamer-armv6l:latest
```
2. Start second container "restreamer2":  
```
docker run -d --name restreamer2 --restart always -e "RESTREAMER_USERNAME=admin" -e "RESTREAMER_PASSWORD=datarhei" -e -p 8090:8080 -v /mnt/restreamer/db2:/restreamer/db --privileged datarhei/restreamer-armv6l:latest
```
3. Open browser and go to your IP with the new port http://yo.ur.device.ip:8090 to configurate your second Datarhei/Restreamer container restreamer2. The default restreamer is reachable on default port 8080.

---
<h1 id="armv7l"></h1>
## Raspberry Pi2 (ARMv7l)    
1. SSH Login to your ARMv7l device
1. Start first image "restreamer" if not already running:  
```
docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=admin" -e "RESTREAMER_PASSWORD=datarhei" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db --privileged datarhei/restreamer-armv7l:latest
```

2. Start second container "restreamer2":  
```
docker run -d --name restreamer2 --restart always -e "RESTREAMER_USERNAME=admin" -e "RESTREAMER_PASSWORD=datarhei" -p 8090:8080 -v /mnt/restreamer/db2:/restreamer/db --privileged datarhei/restreamer-armv7l:latest
```

3. Open browser and go to your IP with the new port http://yo.ur.device.ip:8090 to configurate your second Datarhei/Restreamer container restreamer2. The default restreamer is reachable on default port 8080.

---
## Linux 64 bit
1. SSH Login to your Linux machine 
1. Start first image "restreamer"if not already running:  
```
docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer:latest
```

2. Start second container "restreamer2":  
```
docker run -d --name restreamer2 --restart always -e "RESTREAMER_USERNAME=admin" -e "RESTREAMER_PASSWORD=datarhei" -p 8090:8080 -v /mnt/restreamer/db2:/restreamer/db --privileged datarhei/restreamer:latest
```

3. Open browser and go to your IP with the new port http://yo.ur.device.ip:8090 to configurate your second Datarhei/Restreamer container restreamer2. The default restreamer is reachable on default port 8080.

---
## MacOS and Windows 
1. Open your Docker Toolbox (Kitematic)
2. Open the CLI in Kitematic
2. Start first image "restreamer" if not already running:  
```
docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=admin" -e "RESTREAMER_PASSWORD=datarhei" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db --privileged datarhei/restreamer:latest
```

3. Start second container "restreamer2":  
```
docker run -d --name restreamer2 --restart always -e "RESTREAMER_USERNAME=admin" -e "RESTREAMER_PASSWORD=datarhei" -p 8090:8080 -v /mnt/restreamer/db2:/restreamer/db --privileged datarhei/restreamer:latest
```

4. Open browser and go to your IP with the new port http://yo.ur.device.ip:8090 to configurate your second Datarhei/Restreamer container restreamer2. The default restreamer is reachable on default port 8080.

---
<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>