#!/bin/bash
cpu=$(uname -m | cut -c 1-3)
if [ "${MODE}" == "RASPICAM" ] && [ "$cpu" == "arm" ];
then
    echo "/opt/vc/lib" > /etc/ld.so.conf.d/00-vmcs.conf
    ldconfig
    npm start &
    sleep 20
    /opt/vc/bin/raspivid -t 0 -w 1280 -h 720 -fps 25 -b 500000 -o - | ffmpeg -i - -c copy -f flv rtmp://127.0.0.1:1935/live/raspicam.stream
else
    npm start
fi
