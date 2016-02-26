---
title: guides-youtube-live
---
######[User Guides](../docs/guides-index.html) > YouTube
# Streaming to YouTube-Live over RTMP
To forward your videostream to YouTube-Live, you can use the field "Additional RTMP streaming server". The check box is accessible once the connection to your camera is established.  
<img src="../img/references-external-rtmp-streaming-server.png" width="95%">

##Do the following steps:
1. The premise is that you registered at YouTube to use the live service and of course you have a Datarhei/Restreamer running with a connected video device.

2. Live streaming has to be activated in your YouTube account <a target= "_blank" href="https://support.google.com/youtube/topic/6136989?hl=en&ref_topic=2853712/">(Help)</a>

3. The URL to your live dashboard on YouTube is: <a target= "_blank" href="https://www.youtube.com/live_dashboard">https://www.youtube.com/live_dashboard</a>

4. Find out your RTSP-URL and the uniquie streaming ID like on the image below.  
![YouTube Encoder config](http://datarhei.org/wiki/pic/youtube_rtmp.png)  

5. Combine the RTSP URL an the ID together to get the RTMP URL for the Datarhei/Restreamre output.  rtmp://a.rtmp.youtube.com/live2/[your_channe/IDl] will look like: rtmp://a.rtmp.youtube.com/live2/yourkey-1234.1yyz-4g41-1aq4-abcd.  

6. Next, add the destination address in your Datarhei/Restreamer  
![Datarhei/Restreamer GUI Youtube Livestreaming Link](../img/youtube_rtmp_output.png)

7. Last step: Start process

Once the process has been successfully established, the stream should also be seen at YouTube-Live. It can sometimes take a couple of seconds - be patient.