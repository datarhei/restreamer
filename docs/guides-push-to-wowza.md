---
title: Push to the Wowza-Streaming-Engine
---

## Push to the Wowza-Streaming-Engine

One of the possible applications is to stream a local camera to a streaming service based on the <a target= "_blank" href="https://www.wowza.com/products/streaming-engine/">Wowza-Streaming-Engine</a> without sharing the player on a website. 

<img src="../img/guides-push-to-wowza.png" width="95%">

This can be the scalable <a target= "_blank" href="https://www.wowza.com/products/streaming-cloud/">Wowza Streaming Cloud</a> or a smart company like <a target= "_blank" href="http://www.video-stream-hosting.com/">Video-Stream-Hosting</a> from Germany.

In both cases you have to

1. connect Restreamer with your Camera
2. put the address of your streaming server into the "Additional RTMP streaming server" field   
   `rtmp://[server-address]:1935/[application-name]/live.stream`   
   or if you are using [ModuleSecureURLParams](https://www.wowza.com/forums/content.php?233-How-to-secure-publishing-from-an-RTMP-encoder-that-does-not-support-authentication), take   
   `rtmp://[server-address]/[application-name]/?doPublish=[secure-url-param]/livestream`

That's it!

---
Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too. If you're having a problem while developing, see <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Known Issues</a>.
