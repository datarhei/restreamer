---
title: Easy setup with Docker
---

# Easy setup with Docker
There are three different and ready to start repo images avaiable upon Docker-Hub. Follow examples to run the Restreamer in  5-10 minutes. (other platforms may need different instructions but can be the same): 

* **ARMv6l**
 * [Raspberry Pi 1](#raspberry-pi-1)
* **ARMv7l**
 * [Raspberry Pi 2](#raspberry-pi-2)
 * [Odroid U3](#odroid-u3)
* **Intel/AMD 64bit**
 * [OSX & Windows](#osx-windows)
 * [Linux)](#linux) (e.g. Ubuntu, Debian)

*[Here](hints.html) are some additional hints of running docker containers*

---

### ARMv6l
Requires a armv6 cpu like Raspberry Pi 1

#### Raspberry Pi 1
1. Copy [Hypriot-Image](http://blog.hypriot.com/getting-started-with-docker-on-your-arm-device/) to your SD-Card
3. Run `docker run -d -p 8080:8080 datarhei/restreamer-armv6l:latest`
4. Browse to http://your-device-ip:8080

---

### ARMv7l
Requires ARMv7 CPU like Raspberry Pi 2 or Odroid U3

#### Raspberry Pi 2
1. Copy [Hypriot-Image](http://blog.hypriot.com/getting-started-with-docker-on-your-arm-device/) to your SD-Card
3. Run `docker run -d -p 8080:8080 datarhei/restreamer-armv6l:latest`
4. Browse to http://your-device-ip:8080
5. Default Login: admin / datarhei

#### Odroid U3
1. Download and copy Ubunut 14.04 image to your SD-Card
2. Boot your device
3. Install Docker: `apt-get install -y docker-engine`
4. Run `docker run -d -p 8080:8080 datarhei/restreamer-armv7l:latest`
5. Browse http://your-device-ip:8080
6. Default Login: admin / datarhei

---

### Intel/AMD 64bit
Requires a 64bit Intel or AMD CPU

#### OSX / Windows
1. Install **Docker Toolbox** for [OSX](https://docs.docker.com/engine/installation/mac/) or [Windows](https://docs.docker.com/engine/installation/windows/) of Docker.com
2. Open Kitematic and next the CLI
3. Run `docker run -d -p 8080:8080 datarhei/restreamer:latest`
4. Read our [VirtualBox hints](hints.html#kitematic-virtualbox) to path through external requests on your virtual machine
5. Browse to http://192.168.99.100:8080 (or click on the previewpage upon your kitematc user-interface)
5. Default Login: admin / datarhei

#### Linux
1. Install [Docker Engine](https://docs.docker.com/engine/installation/ubuntulinux/) by docker.com
2. Run `docker run -d -p 8080:8080 datarhei/restreamer:latest`
3. Browse to http://your-device-ip:8080
4. Default Login: admin / datarhei

---

Want to talk to us? Write email open@datarhei.org, go to [Support](../support.html) or choose a nickname and join us on <a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">#datarhei webchat on freenode</a>.

If you're having a weird problem while developing, see [Known Issues](https://github.com/datarhei/restreamer/issues/). 
