---
title: Push to the Wowza-Streaming-Engine
---

One of the possible applications is to stream your local camera to a streaming service based on the <a target= "_blank" href="https://www.wowza.com/products/streaming-engine/">Wowza-Streaming-Engine</a> without sharing the player on a website. 

![Wowza](../img/guides-push-to-wowza.png)

This can be the scalable <a target= "_blank" href="https://www.wowza.com/products/streaming-cloud/">Wowza Streaming Cloud</a> or a smart company like <a target= "_blank" href="http://www.video-stream-hosting.com/">Video-Stream-Hosting</a> from Germany.

In both cases you have to

1. Connect Restreamer with your camera.

2. Put the address of your streaming server into the "External RTMP-Streaming-Server" field.   
   ```
   rtmp://[server-address]/[application-name]/live.stream
   ```

   or if you are using [ModuleSecureURLParams](https://www.wowza.com/forums/content.php?233-How-to-secure-publishing-from-an-RTMP-encoder-that-does-not-support-authentication), use
   ```
   rtmp://[server-address]/[application-name]/?doPublish=[secure-url-param]/livestream
   ```

That's it!
