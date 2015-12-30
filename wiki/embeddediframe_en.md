---
title: embeddediframe_en
---
######[General Instructions](/restreamer/wiki/general_instructions_en.html) > How to embed Snapshot
#How to embed Snapshot
To embed your Snapshot follow the next steps:

1. Open your Datarhei/Restreamer GUI via local IP of your Datarhei/Restreamer device in the webbrowser of you choice  
2. Open the Datarhei/Restreamer player (marked red in screenshot) ![Open your player in GUI](http://datarhei.org/wiki/pic/screenshot_player.jpg) 
3. The HTML Snippet for the video iframe code will be generated under the videoplayer (marked red in screenshot). Add your own dynamic IP and put it on your own website. ![Copy the code and embedd where ever you want](http://datarhei.org/wiki/pic/screenshot_player_iframecode.jpg)  
4. Datarhei/Restreamer generates the html tag with your public internet IP. You can add your own dynamic IP like: `scr="http://123.245.789.123:8008/player.html"` with your [[dynamic IP|dynamic IP]] - for example if your dynamic ip is: foobar.zapto.org the Datarhei/Restreamer link has to be: `http://foobar.zapto.org:8080`. The complete html tag should be like this: `<img src="http://foobar.zapto.org:/player.html" width="800" height="450">`   
5. Change the size over width and height. **Be patient of the correct aspect ratio for a pretty video!***

**Datarhei Hint ☺** ► Do not forget forwarding the TCP Port 8080 and 1934 ([Keyword: Porforwarding](/restreamer/wiki/portforwarding_en.html)) to use the html code successful on your website. How to do portforwarding should be explained in the manual of your router.