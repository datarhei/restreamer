---
title: guides-embed-snapshot
---
#In the following discription we explain how to embed the static snapshot.

1. Open the Datarhei GUI with the local IP of Datarhei in your Browser.  
2. Open the player (marked red in the following Screenshot) ![Open your player in GUI](http://datarhei.org/wiki/pic/screenshot_player.jpg)  
3. You will find the embedd html code under the videoplayer (marked red on the screenshot). Add the HTML with your dynamic IP and put it on your website.![Copy the code and embedd wherever you want](http://datarhei.org/wiki/pic/screenshot_player_snapshot.jpg)  
4. Datarhei generates the html tag with your public internet IP. You can add your own dynamic IP like: `<img src="http://123.456.789.123:8080/images/live.jpg"` with your [[dynamic IP|dynamic IP]] - for example if your dynamic ip is: foobar.zapto.org the Datarhei link has to be: `http://foobar.zapto.org:8080`. The complete html tag should be like this: `<img src="http://foobar.zapto.org:8080/images/live.jpg" width="800" height="450">`   
5. The snapshot will be updated default every 60 seconds. You could change this value in the configuration file of Datarhei. [How to change Snapshotinterval](/restreamer/wiki/modifysnapshotinterval)

**Datarhei Hint ☺** ► Do not forget to forward HTTP Port 8080 and 1934 for a successful Snapshot (Keyword: Porforwarding). Take a look at the manual of your router if you do not know how to do it.  

---

Want to talk to us? Write email open@datarhei.org, go to [Support](../support.html) or choose a nickname and join us on <a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">#datarhei webchat on freenode</a>.

If you're having a weird problem while developing, see [Known Issues](https://github.com/datarhei/small-restreamer-internal/issues/). 
