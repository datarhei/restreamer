---
title: Embed iFrame
---

Here we describe how you embed the player from the Restreamer.

1. Open the Restreamer GUI with the local IP of Restreamer in your browser.  
2. Open the player (marked red in the following Screenshot) ![Open your player in GUI](../img/guides-embed-upon-your-website-1.png)  
3. You will find the embed HTML code under the video player (marked red on the screenshot). Add the HTML with your dynamic IP and put it on your website. ![Copy the code and embedd wherever you want](../img/guides-embed-upon-your-website-2.png)  
4. Restreamer generates the HTML tag with your public internet IP. You can add your own dynamic DNS. Replace the IP in `<iframe src="http://123.456.789.123:8080/player.html" ...>` with your dynamic DNS, for example if your dynamic DNS is: foobar.zapto.org the snapshot URL has to be: `http://foobar.zapto.org:8080`. The HTML tag should look like this: `<iframe src="http://foobar.zapto.org:8080/player.html" ...>`   

Do not forget to [forward the HTTP port](portforwarding.html) 8080 for a successful snapshot. Take a look at the manual of your router if
you do not know how to do it.
{: .notice--info}