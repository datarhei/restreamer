---
title: Installation on Cloud/Hosting
---

# Installation on Cloud/Hosting

If your private ISP has´nt enough bandwith it is obvious to deploy the Restreamer on a Cloud/Hosting-Service. Here is a list of services to deploy a Docker-Image:

**Ready to deploy solutions:**

* <a target= "_blank" href="http://sloppy.io/">sloppy.io</a>
* <a target= "_blank" href="https://tectonic.com/">TECTONIC</a>
* <a target= "_blank" href="https://www.dotcloud.com/">dotCloud</a>
* <a target= "_blank" href="https://www.tutum.co/">tutum</a>

**Cloud solutions:**

* <a target= "_blank" href="https://docs.docker.com/engine/installation/amazon/">Amazon EC2</a>
* <a target= "_blank" href="https://docs.docker.com/engine/installation/google/">Google Cloud Platform]</a>
* <a target= "_blank" href="https://docs.docker.com/engine/installation/softlayer/">IBM SoftLayer</a>
* <a target= "_blank" href="https://docs.docker.com/engine/installation/azure/">Microsoft Azure</a>
* <a target= "_blank" href="https://docs.docker.com/engine/installation/rackspace/">Rackspace Cloud</a>
* <a target= "_blank" href="https://docs.docker.com/engine/installation/joyent/">Joyent’s Triton Elastic Container Service</a>

### Deployment

1. Select and setup your Cloud-/Hosting-Service
2. Start the Restreamer:    
   ```
   # docker run -d --name restreamer --restart always -e "RESTREAMER_USERNAME=YOUR-USERNAME" -e "RESTREAMER_PASSWORD=YOUR-PASSWORD" -p 8080:8080 -v /mnt/restreamer/db:/restreamer/db datarhei/restreamer:latest
   ```
3. Browse to http://your-cloud/hosting-ip:8080

Default login if not set:

* Username: admin
* Password: datarhei

---

## Declaration of the command

* --name restreamer
  you can login into the container by typing "docker exec -it restreamer /bin/bash"
* --restart always   
  the Docker-Daemon is monitoring your container and will start it it again, if it runs into errors. 
* -e "RESTREAMER_USERNAME=..." -e "RESTREAMER_PASSWORD=..."   
  set the login data as enviroment-variable (more [here](references-environment-vars.html#login-security))
* -p 8080:8080   
  bind the device-port 8080 to the Restreamer-port 8080 (you can change it with "-p 31000:8080")
* -v /mnt/restreamer/db:/restreamer/db   
  this save and export the Restreamer-DB on your device-filesystem under /mnt/restreamer

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.