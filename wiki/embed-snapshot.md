---
title: Embed snapshot
---

Here we describe how you embed the static snapshot from the Restreamer.

1. Open the Restreamer GUI with the local IP of Restreamer in your b rowser.
2. Open the player (marked red in the following Screenshot) ![Open your player in GUI](../img/guides-embed-upon-your-website-1.png)
3. You will find the embed HTML code under the video player (marked red on the screenshot). Add the HTML with your dynamic IP and put it on your website.![Copy the code and embedd wherever you want](../img/guides-embed-upon-your-website-3.png)  
4. Restreamer generates the HTML tag with your public internet IP. You can add your own dynamic DNS. Replace the IP in `<img src="http://123.456.789.123:8080/images/live.jpg">` with your dynamic DNS, for example if your dynamic DNS is: foobar.zapto.org the snapshot URL has to be: `http://foobar.zapto.org:8080`. The complete HTML tag should look like this: `<img src="http://foobar.zapto.org:8080/images/live.jpg" width="800" height="450">`   
5. The snapshot will be updated by default every 60 seconds. [How to change the snapshot interval](modify-snapshot-interval.html)

Do not forget to [forward the HTTP port](portforwarding.html) 8080 for a successful snapshot. Take a look at the manual of your router if
you do not know how to do it.
{: .notice--info}