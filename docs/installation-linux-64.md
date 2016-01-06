---
title: Installation under Linux 64bit
---

# Installation Linux 64bit

**Requirements:**

* 64bit Intel or AMD CPU

---

## Installation

1. Download and install Docker-Engine. It's available for: <a target= "_blank" href="https://docs.docker.com/engine/installation/ubuntulinux/">Ubuntu</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/ubuntulinux/">Red Hat Enterprise Linux</a>,  <a target= "_blank" href="https://docs.docker.com/engine/installation/centos/">CentOS</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/fedora/">Fedora</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/debian/">Debian</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/archlinux/">Arch Linux</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/cruxlinux/">CRUX Linux</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/frugalware/">FrugalWare</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/gentoolinux/">Gentoo</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/oracle/">Oracle Linux</a>, <a target= "_blank" href="https://docs.docker.com/engine/installation/SUSE/">openSUSE or SUSE Linux</a>
2. Start the Restreamer:    
   ```
   # docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer:latest
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
* -e `"RESTREAMER_USERNAME=..."` -e `"RESTREAMER_PASSWORD=..."`   
  Set the login data as enviroment-variable (more [here](references-environment-vars.html#login-security))
* -p 8080:8080   
  Bind the device-port 8080 to the Restreamer-port 8080 (you can change it with "-p 31000:8080")
* -v /mnt/restreamer/db:/restreamer/db   
  Save and export the Restreamer-DB on your device-filesystem under /mnt/restreamer

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.
