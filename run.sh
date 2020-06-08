#!/bin/sh
MODE=${MODE:="default"}
RS_MODE=${RS_MODE:=$MODE}

envsubst '$${RS_RTMP_PORT}' < /restreamer/conf/nginx.conf.template > /restreamer/conf/nginx.conf
envsubst '$${RS_RTMP_PORT}' < /restreamer/conf/nginx_ssl.conf.template > /restreamer/conf/nginx_ssl.conf
envsubst '$${RS_RTMP_PORT}' < /restreamer/conf/live.json.template > /restreamer/conf/live.json
envsubst '$${RS_RTMP_PORT}' < /restreamer/conf/live-rpi.json.template > /restreamer/conf/live-rpi.json

# debug reporting
RS_DEBUG=${RS_DEBUG:="false"}
if [ "$RS_DEBUG" = "true" ]; then
    export FFREPORT="file=/restreamer/src/webserver/public/debug/%p-%t.log:level=48"
fi

DEVICE=${DEVICE:="default"}

if [ "$DEVICE" = "raspi" ]; then
    cp ./conf/live-rpi.json ./conf/live.json
fi

echo "/opt/vc/lib" > /etc/ld.so.conf.d/00-vmcs.conf
ldconfig

CPU_TYPE=$(uname -m | cut -c 1-3)
if [ "${RS_MODE}" = "RASPICAM" ] && [ "$CPU_TYPE" = "arm" ]; then
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

    ## options

    RS_RASPICAM_INLINE=${RS_RASPICAM_INLINE:="true"}

    RASPIVID_OPTIONS="--timeout 0 --nopreview"

    # flip the image horizontally
    if [ "$RS_RASPICAM_HFLIP" = "true" ]; then
        RASPIVID_OPTIONS="$RASPIVID_OPTIONS --hflip"
    fi

    # flip the image vertically
    if [ "$RS_RASPICAM_VFLIP" = "true" ]; then
        RASPIVID_OPTIONS="$RASPIVID_OPTIONS --vflip"
    fi

    # inline headers (SPS, PPS) -- https://en.wikipedia.org/wiki/X264
    if [ "$RS_RASPICAM_INLINE" = "true" ]; then
        RASPIVID_OPTIONS="$RASPIVID_OPTIONS --inline"
    fi

    # video stabilization
    if [ "$RS_RASPICAM_STABILIZATION" = "true" ]; then
        RASPIVID_OPTIONS="$RASPIVID_OPTIONS --vstab"
    fi

    ## FPS and GOP

    RASPICAM_FPS=25
    RASPICAM_GOP=50

    if [ -n "$RS_RASPICAM_FPS" ]; then
        RASPICAM_FPS=$RS_RASPICAM_FPS
        RASPICAM_GOP=$((RASPICAM_FPS * 2))
    fi

    if [ -n "$RS_RASPICAM_GOP" ]; then
        RASPICAM_GOP=$RASPICAM_GOP
    fi

    ## bitrate and codec

    RASPICAM_BITRATE=${RS_RASPICAM_BITRATE:=5000000}
    # h264 profile: baseline (no B-frames), main, high
    RASPICAM_H264PROFILE=${RS_RASPICAM_H264PROFILE:="high"}
    # h264 level: 4, 4.1, 4.2 (1080p30, 720p60 and 640 × 480p60/90)
    RASPICAM_H264LEVEL=${RS_RASPICAM_H264LEVEL:=4}
    RASPICAM_CODEC=${RS_RASPICAM_CODEC:="H264"}

    ## image parameter

    RASPICAM_WIDTH=${RS_RASPICAM_WIDTH:=1920}
    RASPICAM_HEIGHT=${RS_RASPICAM_HEIGHT:=1080}
    # -100 to 100
    RASPICAM_SHARPNESS=${RS_RASPICAM_SHARPNESS:=0}
    # -100 to 100
    RASPICAM_CONTRAST=${RS_RASPICAM_CONTRAST:=0}
    # 0 to 100
    RASPICAM_BRIGHTNESS=${RS_RASPICAM_BRIGHTNESS:=50}
    # -100 to 100
    RASPICAM_SATURATION=${RS_RASPICAM_SATURATION:=0}

    ## image quality and effects

    # capture ISO
    RASPICAM_ISO=${RS_RASPICAM_ISO:=0}
    # quantization parameter
    RASPICAM_QP=${RS_RASPICAM_QP:=0}
    # EV compensation, steps of 1/6 stop
    RASPICAM_EV=${RS_RASPICAM_EV:=0}
    # exposure: off,auto,night,nightpreview,backlight,spotlight,sports,snow,beach,verylong,fixedfps,antishake,fireworks
    RASPICAM_EXPOSURE=${RS_RASPICAM_EXPOSURE:="auto"}
    # flicker filter: off,auto,50hz,60hz
    RASPICAM_FLICKER=${RS_RASPICAM_FLICKER:="off"}
    # white balance: off,auto,sun,cloud,shade,tungsten,fluorescent,incandescent,flash,horizon
    RASPICAM_AWB=${RS_RASPICAM_AWB:="auto"}
    # image effects: none,negative,solarise,sketch,denoise,emboss,oilpaint,hatch,gpen,pastel,watercolour,film,blur,saturation,colourswap,washedout,posterise,colourpoint,colourbalance,cartoon
    RASPICAM_IMXFX=${RS_RASPICAM_IMXFX:="none"}
    # metering: average,spot,backlit,matrix
    RASPICAM_METERING=${RS_RASPICAM_METERING:="average"}
    # dynamic range compression: off,low,med,high
    RASPICAM_DRC=${RS_RASPICAM_DRC:="off"}

    ## audio

    RASPICAM_AUDIODEVICE=${RS_RASPICAM_AUDIODEVICE:="0"}
    RASPICAM_AUDIOBITRATE=${RS_RASPICAM_AUDIOBITRATE:="0"}
    RASPICAM_AUDIOCHANNELS=${RS_RASPICAM_AUDIOCHANNELS:="1"}
    RASPICAM_AUDIOLAYOUT=${RS_RASPICAM_AUDIOLAYOUT:="mono"}
    RASPICAM_AUDIOSAMPLING=${RS_RASPICAM_AUDIOSAMPLING:="44100"}

    RASPICAM_AUDIO="-f lavfi -i anullsrc=r=${RASPICAM_AUDIOSAMPLING}:cl=${RASPICAM_AUDIOLAYOUT}"

    if [ "$RS_RASPICAM_AUDIO" = "true" ]; then
        RASPICAM_AUDIO="-thread_queue_size 512 -f alsa -ac ${RASPICAM_AUDIOCHANNELS} -ar ${RASPICAM_AUDIOSAMPLING} -i hw:${RASPICAM_AUDIODEVICE}"
        RASPICAM_AUDIOBITRATE=${RS_RASPICAM_AUDIOBITRATE:="64000"}
    fi

    RASPICAM_AUDIOBITRATE=$(($RASPICAM_AUDIOBITRATE / 1024))

    ## RTMP URL

    RTMP_URL="rtmp://127.0.0.1:${RS_RTMP_PORT}/live/raspicam.stream"

    if [ -n "$RS_TOKEN" ]; then
        RTMP_URL="${RTMP_URL}?token=${RS_TOKEN}"
    fi

    /opt/vc/bin/raspivid \
        $RASPIVID_OPTIONS \
        --width "$RASPICAM_WIDTH" \
        --height "$RASPICAM_HEIGHT" \
        --framerate "$RASPICAM_FPS" \
        --bitrate "$RASPICAM_BITRATE" \
        --intra "$RASPICAM_GOP" \
        --codec "$RASPICAM_CODEC" \
        --profile "$RASPICAM_H264PROFILE" \
        --level "$RASPICAM_H264LEVEL" \
        --sharpness "$RASPICAM_SHARPNESS" \
        --contrast "$RASPICAM_CONTRAST" \
        --brightness "$RASPICAM_BRIGHTNESS" \
        --saturation "$RASPICAM_SATURATION" \
        --ISO "$RASPICAM_ISO" \
        --qp "$RASPICAM_QP" \
        --exposure "$RASPICAM_EXPOSURE" \
        --flicker "$RASPICAM_FLICKER" \
        --awb "$RASPICAM_AWB" \
        --imxfx "$RASPICAM_IMXFX" \
        --metering "$RASPICAM_METERING" \
        --drc "$RASPICAM_DRC" \
        -o - | ffmpeg -i - ${RASPICAM_AUDIO} -map 0:v -map 1:a -codec:v copy -codec:a aac -b:a "${RASPICAM_AUDIOBITRATE}k" -shortest -f flv "${RTMP_URL}" > /dev/null 2>&1
elif [ "${RS_MODE}" = "USBCAM" ]; then
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

    # https://trac.ffmpeg.org/wiki/Capture/Webcam
    # https://www.ffmpeg.org/ffmpeg-devices.html#video4linux2_002c-v4l2
    # https://trac.ffmpeg.org/wiki/Capture/ALSA
    # https://ffmpeg.org/ffmpeg-devices.html#alsa

    USBCAM_VIDEODEVICE=${RS_USBCAM_VIDEODEVICE:="/dev/video"}
    USBCAM_AUDIODEVICE=${RS_USBCAM_AUDIODEVICE:="0"}

    USBCAM_FPS=25
    USBCAM_GOP=50

    if [ -n "$RS_USBCAM_FPS" ]; then
        USBCAM_FPS=$RS_USBCAM_FPS
        USBCAM_GOP=$((USBCAM_FPS * 2))
    fi

    if [ -n "$RS_USBCAM_GOP" ]; then
        USBCAM_GOP=$USBCAM_GOP
    fi

    USBCAM_BITRATE=${RS_USBCAM_BITRATE:=5000000}
    USBCAM_H264PRESET=${RS_USBCAM_H264PRESET:="ultrafast"}
    USBCAM_H264PROFILE=${RS_USBCAM_H264PROFILE:="baseline"}
    USBCAM_WIDTH=${RS_USBCAM_WIDTH:=1280}
    USBCAM_HEIGHT=${RS_USBCAM_HEIGHT:=720}

    USBCAM_BITRATE=$(($USBCAM_BITRATE / 1024))
    USBCAM_BUFFER=$(($USBCAM_BITRATE * 2))

    RTMP_URL="rtmp://127.0.0.1:${RS_RTMP_PORT}/live/usbcam.stream"

    if [ -n "$RS_TOKEN" ]; then
        RTMP_URL="${RTMP_URL}?token=${RS_TOKEN}"
    fi

    USBCAM_VIDEOENCODER="-codec:v libx264 -preset:v ${USBCAM_H264PRESET} -vf format=yuv420p"

    if [ "$DEVICE" = "raspi" ]; then
        USBCAM_VIDEOENCODER="-codec:v h264_omx -profile:v ${USBCAM_H264PROFILE}"
    fi

    USBCAM_AUDIOBITRATE=${RS_USBCAM_AUDIOBITRATE:="0"}
    USBCAM_AUDIOCHANNELS=${RS_USBCAM_AUDIOCHANNELS:="1"}
    USBCAM_AUDIOLAYOUT=${RS_USBCAM_AUDIOLAYOUT:="mono"}
    USBCAM_AUDIOSAMPLING=${RS_USBCAM_AUDIOSAMPLING:="44100"}

    USBCAM_AUDIO="-f lavfi -i anullsrc=r=${USBCAM_AUDIOSAMPLING}:cl=${USBCAM_AUDIOLAYOUT}"

    if [ "$RS_USBCAM_AUDIO" = "true" ]; then
        USBCAM_AUDIO="-thread_queue_size 512 -f alsa -ac ${USBCAM_AUDIOCHANNELS} -ar ${USBCAM_AUDIOSAMPLING} -i hw:${USBCAM_AUDIODEVICE}"
        USBCAM_AUDIOBITRATE=${RS_USBCAM_AUDIOBITRATE:="64000"}
    fi

    USBCAM_AUDIOBITRATE=$(($USBCAM_AUDIOBITRATE / 1024))

    ffmpeg -thread_queue_size 512 -f v4l2 -framerate "$USBCAM_FPS" -video_size "${USBCAM_WIDTH}x${USBCAM_HEIGHT}" -i "${USBCAM_VIDEODEVICE}" ${USBCAM_AUDIO} -map 0:v -map 1:a ${USBCAM_VIDEOENCODER} -r "$USBCAM_FPS" -g "${USBCAM_GOP}" -b:v "${USBCAM_BITRATE}k" -bufsize "${USBCAM_BUFFER}k" -codec:a aac -b:a "${USBCAM_AUDIOBITRATE}k" -shortest -f flv "${RTMP_URL}" > /dev/null 2>&1
else
    npm start
fi
