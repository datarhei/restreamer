---
title: Installation under Linux ARMv6l / ARMv7l
---

# Installation Linux ARMv6l / ARMv7l

**Requirements:**

* ARMv6l or ARMv7l architecture

---
## ARMv6l
Verified with <a target= "_blank" href="https://www.raspberrypi.org/products/model-b-plus/">Raspberry Pi 1</a>:

### Raspberry Pi 1

1. Copy the <a target= "_blank" href="http://blog.hypriot.com/getting-started-with-docker-on-your-arm-device/">Hypriot-Image</a> to your SD-Card or download and install the <a target= "_blank" href="http://blog.hypriot.com/downloads/#hypriot-docker-debian-dackages-for-raspberry-pi">Hypriot Docker Debian Packages for Raspberry Pi</a>  
2. login to your ARMv6l-Board and start the Restreamer with:    
   ```
   # docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer-armv6l:latest
   ```
3. Browse to http://your-device-ip:8080

The default login (more [here](references-environment-vars.html#login-security)) are:

* Username: admin
* Password: datarhei

---

## ARMv7l
Tested with <a target= "_blank" href="https://www.raspberrypi.org/products/raspberry-pi-2-model-b/">Raspberry Pi 2</a> and <a target= "_blank" href="http://www.hardkernel.com/main/products/prdt_info.php?g_code=g138745696275">Odroid U3</a>

### Raspberry Pi 2

1. Copy the <a target= "_blank" href="http://blog.hypriot.com/getting-started-with-docker-on-your-arm-device/">Hypriot-Image</a> to your SD-Card or download and install the <a target= "_blank" href="http://blog.hypriot.com/downloads/#hypriot-docker-debian-dackages-for-raspberry-pi">Hypriot Docker Debian Packages for Raspberry Pi</a>  
2. start the Restreamer with:    
   ```
   # docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer-armv7l:latest
   ```
3. Browse to http://your-device-ip:8080

The default login (more [here](references-environment-vars.html#login-security)) are:

* Username: admin
* Password: datarhei

### Odroid U3

1. Download and copy Ubunut 14.04 image to your SD-Card
2. Install Docker-Engine (tested with kernel 3.8.13.30):    
   ```
   apt-get update && apt-get install -y docker-engine
   ```
3. Start the Restreamer with:    
   ```
   # docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer-armv7l:latest
   ```
4. Browse to http://your-device-ip:8080

The default login (more [here](references-environment-vars.html#login-security)) are:

* Username: admin
* Password: datarhei

---

## Declaration of the command

* --name restreamer
  you can login into the container by typing "docker exec -it restreamer /bin/bash"
* --restart always   
  the Docker-Daemon is monitoring your container and will start it it again, if it runs into errors. 
* -e `"RESTREAMER_USERNAME=..."` -e `"RESTREAMER_PASSWORD=..."`   
  set the login data as enviroment-variable (more [here](references-environment-vars.html#login-security))
* -p 8080:8080   
  bind the device-port 8080 to the Restreamer-port 8080 (you can change it with "-p 31000:8080")
* -v /mnt/restreamer/db:/restreamer/db   
  this save and export the Restreamer-DB on your device-filesystem under /mnt/restreamer

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.
