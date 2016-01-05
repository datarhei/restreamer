---
title: Configuration
---

# Configuration

The base configuration file is located on "/config/config.json". Please don't change its settings if you don't understand the architecture and the software components like NGINX and FFmpeg..

The user configurations are available at [Enviroment-Variables](references-environment-vars.html).

---

#### name

```json
{
    "name": "live"
}
```
If you want to add more than one config.json, you have to change this name to "dev" or something you like and start the application with the env. "MODE=dev" to select the other one.

#### jsondb

```json
{
    "jsondb": "db/v1"
}
```

This sets the path to the Jsondb. More about the Jsondb  [here](https://www.npmjs.com/package/node-jsondb)

#### auth

```json
{
    "auth": {
        "username": "admin",   
        "password": "datarhei" 
    }
}
```
Just a fallback for the login data if no env. is set. More about the env. [here](references-environment-vars.html#login-security)

#### ffmpeg options

```json
{
    "ffmpeg": {
        "options": {
            "native_h264":[  
                "-c copy",
                "-f flv"
            ],
            "native_h264_soundless_aac":[  
                "-ar 44100",
                "-ac 2",
                "-acodec pcm_s16le",
                "-f s16le",
                "-ac 2",
                "-i /dev/zero",
                "-c:v copy",
                "-acodec aac",
                "-ab 128k",
                "-f flv"
            ],
            "snapshot": "-vframes 1"  
        }
    }
}
```

If FFmpeg is initiated, Restreamer will try to check if your video stream has an audio channel or not. It will choose "nativeh264" if your stream has an audio and "nativeh264soundlessaac" if not. This is required if you want to stream a video to YouTube-Live without native audio. For the creation of a snapshot we select the "snapshot" options.
        
#### ffmpeg monitor

```json
{
    "ffmpeg": {
        "monitor": {
            "restart_wait": "6000",  
            "retries": 10  
        }
    }
}
```

If a FFmpeg-process dies, the monitor will try to restart it. "restart_wait" sets the time the monitor has to wait before retrying and "retries" is the count for the monitor to start the process again.

#### nginx

```json
{
    "nginx": {
        "exec": "/usr/local/nginx/sbin/nginx -c /restreamer/conf/nginx.conf",  
        "streaming": {
            "ip": "127.0.0.1",  
            "rtmp_port": "1935",  
            "rtmp_path": "/live/",  
            "hls_port": "8080", 
            "hls_path": "/hls/" 
        }
    }
}
```

These options are required if you want to use an external streaming server from our docker container. Warning: the player and a lot of other processes are currently static on the local NGINX. Actually it is not recommended to change it!  

"exec" is the command the NodeJS-Application of Restreamer will execute on startup to launch the NGINX-instance.   
"ip" sets the IP to the NGINX-RTMP server   
"rtmp_port" set the rtmp port   
"rtmp_path" set the path to the rtmp-application   
"hls_port" set the hls port   
"hls_path" set the path to the hls-application   

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.