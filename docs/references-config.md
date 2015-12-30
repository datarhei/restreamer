---
title: Config
---

# Config

```json
{
    // default conf
    "name": "live", 
    // path to jsondb file
    "jsondb": "db/v1",  
    "auth": {
        // user-interface username
        "username": "admin",   
        // user-interface password
        "password": "datarhei" 
    },
    "ffmpeg": {
        "options": {
            // ffmpeg-options for streams with audio
            "native_h264":[  
                "-c copy",
                "-f flv"
            ],
            // ffmpeg-options for streams without audio (req. for youtube-live)
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
            // ffmpeg-options for the snapshot
            "snapshot": "-vframes 1"  
        },
        "monitor": {
            // time to wait before retry (in ms)
            "restart_wait": "6000",  
            // count for retry
            "retries": 10  
        }
    },
    "nginx": {
        // path to nginx-bin
        "exec": "/usr/local/nginx/sbin/nginx -c /restreamer/conf/nginx.conf",  
        "streaming": {
            // nginx ip
            "ip": "127.0.0.1",  
            // nginx rtmp port
            "rtmp_port": "1935",  
            // nginx rtmp path
            "rtmp_path": "/live/",  
            // nginx hls port
            "hls_port": "8080", 
            // nginx hls path
            "hls_path": "/hls/" 
        }
    }
}
```

---

Want to talk to us? Write email open@datarhei.org, go to [Support](../support.html) or choose a nickname and join us on <a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">#datarhei webchat on freenode</a>.

If you're having a weird problem while developing, see [Known Issues](https://github.com/datarhei/small-restreamer-internal/issues/). 
