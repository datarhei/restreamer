#!/bin/bash
CPU_TYPE=$(uname -m | cut -c 1-3)
if [ "${MODE}" == "RASPICAM" ] && [ "$CPU_TYPE" == "arm" ];
then
    echo "/opt/vc/lib" > /etc/ld.so.conf.d/00-vmcs.conf
    ldconfig
    npm start &
    NGINX_RUNNING=0
    until [ "$NGINX_RUNNING" = "1" ]; do
        if pgrep "nginx" > /dev/null
        then
            NGINX_RUNNING=1
        else
            NGINX_RUNNING=0
            sleep 5
        fi
    done
    /opt/vc/bin/raspivid -t 0 -w 1280 -h 720 -fps 25 -b 500000 -o - | ffmpeg -i - -f lavfi -i aevalsrc=0 -vcodec copy -acodec aac -strict experimental -map 0:0 -map 1:0 -shortest -flags +global_header -f flv rtmp://127.0.0.1:1935/live/raspicam.stream > /dev/null 2>&1
elif [ "${MODE}" == "USBCAM" ];
then
    apt-get update && apt-get install -y v4l-utils libv4l-0
    npm start &
    NGINX_RUNNING=0
    until [ "$NGINX_RUNNING" = "1" ]; do
        if pgrep "nginx" > /dev/null
        then
            NGINX_RUNNING=1
        else
            NGINX_RUNNING=0
            sleep 5
        fi
    done
    ffmpeg -f v4l2 -r 25 -s 1280x720 -i /dev/video0 -f lavfi -i aevalsrc=0 -vcodec copy -acodec aac -strict experimental -map 0:0 -map 1:0 -shortest -flags +global_header -f flv rtmp://127.0.0.1:1935/live/usb.stream > /dev/null 2>&1
else
    npm start
fi
