<h1 align="center">Restreamer</h1>
<h3 align="center">Smart free video streaming.</h3>
<p align="center"><a href="https://demo.datarhei.com/ui">Try live demo</a><br />Username: admin<br />Password: demo</p>
<p align="center"><a href="https://docs.datarhei.com/restreamer">Documentation</a></p>

<p align="center">
  <a href="https://datarhei.com">
    <img src="https://raw.githubusercontent.com/datarhei/restreamer/2.x/readme.png" alt="Restreamer Interface" />
  </a>
</p>

<p align="center">Self-hosting solution to stream live to your website and publish to many like YouTube-Live, Twitter, Twitch, Vimeo, and other platforms or services. Our Docker-Image is easy to install and runs on Linux environments (MacOS/Windows by Docker Desktop). Moreover, combine the Restreamer with single-board computers like Raspberry Pi or GPU powered systems for Video-Encoding.</p>
<br />
<hr />

## Features

- Simplified User-Interface
- REST-API (JSON) and 100% Swagger documented
- Multiple A/V Inputs, Outputs, protocols, and codecs
- Option to mux a separate audio channel to the video
- Raspberry Pi (MMAL/OMX), Nvidia Cuda, Intel VAAPI support
- Support for Hardware- and Virtual-Devices
- FFmpeg Video-Processing (as native as possible)
- Build-in VideoJS-Player for your Website
- Configurable publication website for streaming without player embedding
- HTTP/S- (HLS) and RTMP/S-Streaming Server
- Automatic Let's Encrypt HTTPS certification
- Content license with Creative Commons
- Easy wizard configuration
- Viewer/Bandwidth Monitoring and limiting
- Resource Monitoring (optionally by Prom-Metrics)
- Server- and Process-Logging
- GDPR compliant without third-party providers and does not save audience data

## Quick setup

**AMD64/ARMv7/ARM64:**
```sh
docker run -d --restart=always -name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   -p 8080:8080 -p 8181:8181 -p 1935:1935 datarhei/restreamer:latest
```

*`--privileged` just for local devces like soundcard or hdmi-out.*

**ARMv7 Raspberry Pi:**
```sh
docker run -d --restart=always -name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   --privileged 
   -p 8080:8080 -p 8181:8181 -p 1935:1935 datarhei/restreamer:rpi-lastest
```

**AMD64 Nvidia Cuda:**
```sh
docker run -d --restart=always -name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   --runtime=nvidia 
   -p 8080:8080 -p 8181:8181 -p 1935:1935 datarhei/restreamer:cuda-lastest
```

**AMD64 Intel VAAPI:**
```sh
docker run -d --restart=always -name restreamer \
   -v /opt/restreamer/config:/core/config -v /opt/restreamer/data:/core/data \
   --privileged \
   -p 8080:8080 -p 8181:8181 -p 1935:1935 datarhei/restreamer:vaapi-lastest
```


*For external access, port forwarding from 80/TCP to 8080/TCP and 443/TCP to 8181/TCP to the Restreamer's internal IP address may need to be set up.*

## Documentation

Documentatio is available on [docs.datarhei.com/restreamer](https://docs.datarhei.com/restreamer). We give many pieces of information, from setting up a camera, embedding your player upon your website, and streaming to services like, e.g., YouTube-Live, and many more.

- Quick start
- Installation
- Manual
- Guides
- Developer

## Development

**For the Restreamer interface:**

```
$ git clone github.com/datarhei/restreamer
$ cd restreamer
$ yarn install
$ npm run start
```

**To add/fix translations:**
Locales are located in `src/locals`
```
$ npm run i18n-extract:clean
$ npm run i18n-compile
```

Learn more about the datarhei Core in our Repository. 

## Community support

For general help using Restreamer, please refer to the official [documentation](https://docs.datarhei.com/restreamer). For additional support, you can use Github to ask a question (Bug reports, Contributions, Features).

## License
See the [LICENSE](./LICENSE) file for licensing information.
