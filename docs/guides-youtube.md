---
title: YouTube
---

## Setup Restreamer

In order to be able to push your stream to YouTube, its audio has to be either encoded in `AAC` or `MP3`. If this is not already
the case you can force the encoding of the audio stream by setting the environment variable `RS_AUDIO` to either `aac` or `mp3`.

If you do not want to push the audio of your stream to YouTube, set the environment variable `RS_AUDIO` to `silence`. This will
replace the audio track with silence, encoded in `AAC`.

Check out the [environment variables reference](references-environment-vars.html) for how to set the value value for `RS_AUDIO`.

Read more about the [stream requirements](https://support.google.com/youtube/answer/2853702).

## Setup YouTube

- Log in into your YouTube account
- Make sure that your YouTube channel is enabled for live streaming. [Read more](https://support.google.com/youtube/answer/2474026).
- Go to your [live dashboard](https://www.youtube.com/live_dashboard)

Scroll down and you'll find the YouTube RTMP-URL and your personal unique streaming ID after clicking on "reveal".
![YouTube RTMP-URL](../img/guides-youtube-key-reveal.png)

Your RTMP-URL that you need for Restreamer will look like `rtmp://a.rtmp.youtube.com/live2/0dqe-1cgd-438p-8abw`.

## Push to YouTube

To push your video stream to YouTube, you have to check the field "External RTMP-Streaming-Server" and enter the RTMP-URL from above into the field ...
![YouTube Streaming](../img/guides-rtsp-url-streaming-external-youtube.png)

... and press "Start":

![YouTube Streaming](../img/guides-rtsp-url-streaming-external-youtube-streaming.png)

Once the process has been successfully initiated, the stream should also be seen at YouTube-Live.
Please be patient, this may take a few of seconds.

![Perfect Stream Health](../img/guides-youtube-streamhealth-perfect.png)

## Stream Health

Depending on the bitrate of your stream, YouTube may complain a bit (changing the color from <span style="color:green">green</span> to <span style="color:orange">orange</span>), but this is OK.

![OK Stream Health](../img/guides-youtube-streamhealth-ok.png) ![Stream Health](../img/guides-youtube-streamhealth.png)

On the managmenet interface of your [IP camera](guides-ipcam-rtsp.html), you often have the possibility (depending on the manufacturer) to adjust
the bitrate of the output stream. Please read more about the [stream requirements](https://support.google.com/youtube/answer/2853702).
