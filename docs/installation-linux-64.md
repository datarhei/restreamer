---
title: Installation Linux 64bit
---
###### [Installation](../docs/installation-index.html) > Installation Linux 64bit
# Installation Linux 64bit

**Requirements:**

* 64bit Intel or AMD CPU

---

## Installation

1. Download and install Docker-Engine. It's available for: <a target= "_blank" href="https://docs.docker.com/engine/installation/ubuntulinux/">Ubuntu</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/ubuntulinux/">Red Hat Enterprise Linux</a>,  <a target= "_blank" href="https://docs.docker.com/engine/installation/centos/">CentOS</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/fedora/">Fedora</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/debian/">Debian</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/archlinux/">Arch Linux</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/cruxlinux/">CRUX Linux</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/frugalware/">FrugalWare</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/gentoolinux/">Gentoo</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/oracle/">Oracle Linux</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/SUSE/">openSUSE or SUSE Linux</a>
2. Start the Restreamer:    
   ```
   # docker run -d --name restreamer --restart always -e "RS_USERNAME=YOUR-USERNAME" -e "RS_PASSWORD=YOUR-PASSWORD" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer:latest
   ```
3. Browse to http://your-device-ip:8080

The default login (more [here](references-environment-vars.html#login-security)) are:

* Username: admin
* Password: datarhei

---

## Declaration of the command

* --name restreamer
  You can login into the container by typing "docker exec -it restreamer /bin/bash"
* --restart always   
  The Docker-Daemon is monitoring your container and will start it it again, if it runs into errors. 
* -e `"RS_USERNAME=..."` -e `"RS_PASSWORD=..."`   
  Set the login data as enviroment-variable (more [here](references-environment-vars.html#login-security))
* -p 8080:8080   
  Bind the device-port 8080 to the Restreamer-port 8080 (you can change it with "-p 31000:8080")
* -v /mnt/restreamer/db:/restreamer/db   
  Save and export the Restreamer-DB on your device-filesystem under /mnt/restreamer

---