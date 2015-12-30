---
title: Environment Variables
---

# Enviroment Variables

---

```sh
[INFO]      _       _             _           _  
[INFO]   __| | __ _| |_ __ _ _ __| |___   ___(_) 
[INFO]  / _  |/ _  | __/ _  |  __|  _  |/  _ | | 
[INFO] | (_| | (_| | || (_| | |  | | | |  __/| | 
[INFO] |_____|_____|_||_____|_|  |_| |_|____||_| 
[INFO]  
[INFO] ReStreamer v0.0.1 
[INFO]  
[INFO] ENVIRONMENTS 
[INFO] ENV "NODEJS_PORT=3000"
[INFO] ENV "LOGGER_LEVEL=3"
[INFO] ENV "TIMEZONE=Europe/Berlin"
[INFO] ENV "SNAPSHOT_REFRESH_INTERVAL=60000"
[INFO] ENV "CREATE_HEAPDUMPS=false"
[INFO]
```

---

**Docker example:**

```sh
docker run ...
    -e "LOGGER_LEVEL=4" \
    -e "SNAPSHOT_REFRESH_INTERVAL=10000" \
        ...
```

---

#### NODEJS_PORT

Webserver port of application 

---

#### LOGGER_LEVEL

Logger level to defined, what should be logged

---

#### TIMEZONE

Set the timezone for logging

---

#### SNAPSHOT_REFRESH_INTERVAL

Interval to create a new Snapshot in milliseconds

---

#### CREATE_HEAPDUMPS (dev.)

Create Heapdumps of application (needs g++, make and python to run, please install heapdump with 'npm install heapdump' afterwords)  

---

Want to talk to us? Write email open@datarhei.org, go to [Support](../support.html) or choose a nickname and join us on <a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">#datarhei webchat on freenode</a>.

If you're having a weird problem while developing, see [Known Issues](https://github.com/datarhei/small-restreamer-internal/issues/). 

