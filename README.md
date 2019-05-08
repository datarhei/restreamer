# Restreamer

Datarhei/Restreamer offers smart free video streaming. Stream H.264 video of IP cameras live to your website. Upload your live video on [YouTube-Live](https://www.youtube.com/), [IBM Watson](https://video.ibm.com/), [Twitch](https://www.twitch.tv/), [Vimeo](https://livestream.com/) or any other streaming solutions e.g. [Wowza-Streaming-Engine](https://www.wowza.com/). Our [Docker-Image](https://hub.docker.com/search/?q=restreamer&page=1&isAutomated=0&isOfficial=0&starCount=0&pullCount=0) is easy to install and runs on Linux, MacOS and Windows. Datarhei/Restreamer can be perfectly combined with single-board computers like [Raspberry Pi](https://www.raspberrypi.org/) and [Odroid](http://www.hardkernel.com/main/main.php). It is free (licensed under Apache 2.0) and you can use it for any purpose, private or commercial.  

## Features

- User-Interface including login-security
- JSON / HTTP-API
- <a target= "_blank" href="http://ffmpeg.org/">FFmpeg</a> streaming/encoding the video/camera-stream, creating snapshots or pushing to a external streaming-endpoint
- <a target= "_blank" href="http://nginx.org/">NGINX</a> incl. <a target= "_blank" href="https://github.com/sergey-dryabzhinsky/nginx-rtmp-module">RTMP-Module</a> as streaming-backend and hls server
- <a target= "_blank" href="https://github.com/clappr/clappr">Clappr-Player</a> to embed your stream on your website
- <a target= "_blank" href="https://www.docker.com/">Docker</a> and <a target= "_blank" href="https://kitematic.com/">Kitematic (Docker-Toolbox)</a> optimizations and very easy installation

## Documentation

Documentation is available on [Datarhei/Restreamer GitHub pages](https://datarhei.github.io/restreamer/).
We give you a lot of of informations from setting up a camera, embedding your player upon your website and streaming to services like e.g. YouTube-Live, Ustream and Livestream.com and many more things. 

More additional informations about streaming, cameras and so on you can find in our [Wiki](https://datarhei.github.com/restreamer/wiki). 

## Development

#### Building your own Docker-Image:

```sh
$ git clone https://github.com/datarhei/restreamer
$ docker build -t restreamer .
```

*Required Docker version >= 17.05*

## Help / Bugs

If you have problems or found a bug feel free to create a new issue upon the <a target= "_blank" href="https://github.com/datarhei/restreamer/issues">Github issue management</a>.

Want to talk to us? Write an email to <a href="mailto:open@datarhei.org?subject=Datarhei/Restreamer">open@datarhei.org</a> or ask a question in our (<a target= "_blank" href="https://groups.google.com/forum/#!forum/datarhei">Forum</a>) on Google Groups.

## Authors

The Datarhei/Restreamer was created by [Julius Eitzen](https://github.com/jeitzen), [Sven Erbeck](https://github.com/svenerbeck), [Christoph Johannsdotter](https://github.com/christophjohannsdotter) and [Jan Stabenow](https://github.com/jstabenow).   

Special thanks for supporting this project [Andrew Shulgin](https://github.com/andrew-shulgin).

## Copyright

Code released under the [Apache license](LICENSE). Images are copyrighted by datarhei.org
