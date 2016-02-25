<<<<<<< HEAD
# Restreamer Docs and Website

## Setup

1. Clone the Restreamer repository

2. Change into the "docs" directory where docs live
=======
#Restreamer

Datarhei/Restreamer offers smart free video streaming in real time. Stream H.264 video of IP cameras live to your website. Upload your live video on [YouTube-Live](https://www.youtube.com/), [Ustream](http://www.ustream.tv/), [Twitch](http://www.twitch.tv/), [Livestream.com](http://livestream.com/) or any other streaming solutions e.g. [Wowza-Streaming-Engine](https://www.wowza.com/). Our [Docker-Image](https://hub.docker.com/search/?q=restreamer&page=1&isAutomated=0&isOfficial=0&starCount=0&pullCount=0) is easy to install and runs on Linux, MacOS and Windows. Datarhei/Restreamer can be perfectly combined with single-board computers like [Raspberry Pi](https://www.raspberrypi.org/) and [Odroid](http://www.hardkernel.com/main/main.php). It is free (licensed under Apache 2.0) and you can use it for any purpose, private or commercial.  

![amd64](http://dockeri.co/image/datarhei/restreamer)
![armv6l](http://dockeri.co/image/datarhei/restreamer-armv6l)
![armv6l](http://dockeri.co/image/datarhei/restreamer-armv7l)

##Features
- User-Interface including login-security
- JSON / HTTP-API
- <a target= "_blank" href="http://ffmpeg.org/">FFmpeg</a> streaming/encoding the video/camera-stream, creating snapshots or pushing to a external streaming-endpoint
- <a target= "_blank" href="http://nginx.org/">NGINX</a> incl. <a target= "_blank" href="https://github.com/arut/nginx-rtmp-module">RTMP-Module</a> as streaming-backend and hls server
- <a target= "_blank" href="https://github.com/clappr/clappr">Clappr-Player</a> to embed your stream on your website
- <a target= "_blank" href="https://www.docker.com/">Docker</a> and <a target= "_blank" href="https://kitematic.com/">Kitematic (Docker-Toolbox)</a> optimizations and very easy installation

##Roadmap
- RC6 (tba)

##Documentation
Documentation is available on [Datarhei/Restreamer GitHub pages](https://datarhei.github.io/restreamer/).
We give you a lot of of informations from setting up a camera, embedding your player upon your website and streaming to services like e.g. YouTube-Live, Ustream and Livestream.com and many more things. 
>>>>>>> master

        $ cd docs

<<<<<<< HEAD
3. Build the Dockerfile

        $ docker build -t rs-docs .

4. Start the web server

        $ docker run -d -p 4000:4000 -v $PWD:/restreamer rs-docs

5. Visit the site at
   http://docker-host-ip:4000/restreamer/

## Deployment

1. Clone a separate copy of the Restreamer repo as a sibling of your normal
   Restreamer project directory and name it "restreamer-gh-pages".

        $ git clone git@github.com:datarhei/restreamer.git restreamer-gh-pages

2. Check out the "gh-pages" branch.

        $ cd /path/to/restreamer-gh-pages
        $ git checkout gh-pages

3. Copy the contents of the "docs" directory in master to the root of your
   restreamer-gh-pages directory.

        $ cd /path/to/restreamer
        $ cp -r docs/** ../restreamer-gh-pages

4. Change to the small-restreamer-gh-pages directory, commit, and push the changes

        $ cd /path/to/restreamer-gh-pages
        $ git commit . -m "Syncing docs with master branch"
        $ git push
=======
##Help / Bugs
If you have problems or found a bug feel free to create a new issue upon the <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Github issue management</a>.

Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a>, go to [Support](../support.html) or choose a nickname speak to us in IRC: <a href="irc://irc.freenode.net#piwik">irc.freenode.net/#datarhei</a> (<a target= "_blank" href="https://webchat.freenode.net/?channels=datarhei">webchat</a>). You could ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups, too.

##Authors
The Datarhei/Restreamer was created by [Julius Eitzen](https://github.com/jeitzen), [Sven Erbeck](https://github.com/svenerbeck), [Christoph Johannsdotter](https://github.com/christophjohannsdotter) and [Jan Stabenow](https://github.com/jstabenow).

##Copyright
Code released under the [Apache license](LICENSE). Images are copyrighted by datarhei.org
>>>>>>> master
