<h1 align="center">Restreamer</h1>
<h3 align="center">A really nice and free alternative for handling live streams.</h3>
<p align="center">
<a href="https://github.com/datarhei/restreamer/blob/2.x/LICENSE" target="_blank"><img src="https://img.shields.io/github/license/datarhei/restreamer" alt="License" /></a>
<a href="https://github.com/datarhei/restreamer/releases" target="_blank"><img src="https://img.shields.io/github/v/release/datarhei/restreamer?include_prereleases" alt="License" /></a>
<a href="https://hub.docker.com/r/datarhei/restreamer" target="_blank"><img src="https://img.shields.io/docker/pulls/datarhei/restreamer" alt="Docker pulls" /></a>
<a href="https://docs.datarhei.com/restreamer/getting-started/quick-start" target="_blank"><img src="https://img.shields.io/badge/documentation-get%20started-green" alt="Documentation" /></a>
</p>
<p align="center"><a href="https://demo.datarhei.com/ui" target="_blank">Try live demo</a><br />
<a href="https://demo.datarhei.com/ui" target="_blank"><img src="https://img.shields.io/badge/username-admin-blue" alt="demo username" /></a>
<a href="https://demo.datarhei.com/ui" target="_blank"><img src="https://img.shields.io/badge/password-demo-blue" alt="demo password" /></a>

<p align="center">
  <a href="https://datarhei.com">
    <img src="https://github.com/datarhei/restreamer/blob/2.x/readme-promo.gif" alt="Restreamer Promo Video" />
  </a>
</p>

<p align="center">Self-hosting solution to stream live to your website and publish to many like YouTube-Live, Twitter, Twitch, Vimeo, and other platforms or services. Our Docker-Image is easy to install and runs on Linux environments (MacOS/Windows by Docker Desktop). Moreover, combine the Restreamer with single-board computers like Raspberry Pi or GPU powered systems for Video-Encoding.</p>
<br />
<hr />

## Features

- Simplified User-Interface
- Easy wizard configuration
- Multiple audio/video inputs, outputs, protocols, and codecs
- ReStreaming to platforms such as YouTube-Live, to software such as Wowza Media Server, and others based on protocols such as RTMP, SRT, ...
- Option to mux a separate audio channel to the video
- Build-in VideoJS-Player for your Website
- Configurable publication website for streaming without player embedding
- Content license with Creative Commons
- HTTP/S- (HLS), RTMP/S- & SRT-Streaming Server
- Automatic Let's Encrypt HTTPS certification
- Viewer/Bandwidth Monitoring and limiting
- Raspberry Pi (MMAL/OMX), Nvidia Cuda, Intel VAAPI support
- Support for Hardware- and Virtual-Devices
- FFmpeg Video-Processing (as native as possible)
- REST-API (JSON) and 100% Swagger documented
- Resource Monitoring (optionally by Prom-Metrics)
- Server- and Process-Logging
- GDPR compliant without third-party providers and does not save audience data

## Quick setup

### AMD64/ARMv7/ARM64:
```sh
docker run -d --restart=always --name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   -p 8080:8080 -p 8181:8181 \
   -p 1935:1935 -p 1936:1936 \
   -p 6000:6000/udp \
   datarhei/restreamer:latest
```

*`--privileged` just for local devices like usb cameras.*    
*Try `--security-opt seccomp=unconfined` if no network source can be reached.*

### ARMv7/ARM64 Raspberry Pi:
```sh
docker run -d --restart=always --name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   --privileged \
   -p 8080:8080 -p 8181:8181 \
   -p 1935:1935 -p 1936:1936 \
   -p 6000:6000/udp \
   datarhei/restreamer:rpi-latest
```

*`--privileged` just for local devices like usb cameras.*    
*Try `--security-opt seccomp=unconfined` if no network source can be reached.*

### AMD64 Nvidia Cuda:
```sh
docker run -d --restart=always --name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   --runtime=nvidia --privileged \
   -p 8080:8080 -p 8181:8181 \
   -p 1935:1935 -p 1936:1936 \
   -p 6000:6000/udp \
   datarhei/restreamer:cuda-latest
```

*`--privileged` just for local devices like usb cameras.*    
*Try `--security-opt seccomp=unconfined` if no network source can be reached.*

### AMD64 Intel VAAPI:
```sh
docker run -d --restart=always --name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   -v /dev/dri:/dev/dri --privileged \
   -p 8080:8080 -p 8181:8181 \
   -p 1935:1935 -p 1936:1936 \
   -p 6000:6000/udp \
   datarhei/restreamer:vaapi-latest
```

*`--privileged` just for local devices like usb cameras.*    
*Try `--security-opt seccomp=unconfined` if no network source can be reached.*

*For external access (http/s, rtmp/s, srt), port forwarding from your Internet-Router to the Restreamer's internal IP address may need to be set up.*

## Documentation

Documentation is available on [docs.datarhei.com/restreamer](https://docs.datarhei.com/restreamer). We give many pieces of information, from setting up a camera, embedding your player upon your website, and streaming to services like, e.g., YouTube-Live, and many more.

- [Quick start](https://docs.datarhei.com/restreamer/getting-started/quick-start)
- [Installation](https://docs.datarhei.com/restreamer/installing/minimum-requirements)
- [Manual](https://docs.datarhei.com/restreamer/knowledge-base/manual)
- [Guides](https://docs.datarhei.com/restreamer/knowledge-base/user-guides)

## Development

### Create a custom image (bundle):

#### [Restreamer FFmpeg](https://github.com/datarhei/ffmpeg):
```
$ git clone github.com/datarhei/ffmpeg
$ cd ffmpeg
$ docker build -f Dockerfile.alpine -t myffmpeg .
```

#### [Restreamer backend](https://github.com/datarhei/core) (Golang):

```
$ git clone github.com/datarhei/core
$ cd core
$ docker build -t mycore .
```

#### [Restreamer interface](https://github.com/datarhei/restreamer-ui) (React):
```
$ git clone github.com/datarhei/restreamer-ui
$ cd restreamer-ui
$ docker build -t myrsui .
```

#### Restreamer bundle:
```
$ git clone github.com/datarhei/restreamer
$ cd restreamer
$ docker build --build-arg FFMPEG_IMAGE=myffmpeg --build-arg CORE_IMAGE=mycore --build-arg RESTREAMER_UI_IMAGE=myrsui -t myrestreamer .
$ docker run -it --rm -p 8080:8080 myrestreamer
```

### To add/fix translations in the [Restreamer interface](https://github.com/datarhei/restreamer-ui):

The Restreamer interface is currently translated in different languages, such as German, French, Italian, Spanish, and more. If you find errors in the translations or have better suggestions for some sentences, you can become a translation contributor on [poeditor.com](https://poeditor.com/join/project/ogATl3F48K).
There you can also start a translation into a language that is not yet available in the Restreamer interface.

Contribute to the translations on: [https://poeditor.com/join/project/ogATl3F48K](https://poeditor.com/join/project/ogATl3F48K)

## Community support

For general help using Restreamer, please refer to the official [documentation](https://docs.datarhei.com/restreamer). For additional support, you can use Github to ask a question (Bug reports, Contributions, Features).

## License
See the [LICENSE](./LICENSE) file for licensing information.

## Business inquiries
**We provide support for commercial requirements with professional support, agile software development, and consulting.** If you have a commercial request, be it a bug or a feature enhancement, please contact us directly at support@datarhei.com.
