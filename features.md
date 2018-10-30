---
title: Features
sidebar:
  nav: "features"
---

Restreamer is free open source live video streaming on your website without a streaming provider.

## Core Features

### <i class="fas fa-check"></i> Live Video Streaming
Restreamer is an out-of-the-box free video streaming server and ready to stream directly on your website. [Read more](docs/guides-embedding.html).

### <i class="fas fa-check"></i> H.264
Restreamer supports end-to-end H.264 streaming from the input to the target without quality loss.

### <i class="fas fa-check"></i> HLS Video
Restreamer does not require Adobe Flash. It uses HLS HTTP Streaming which is the standart now. [Read more](wiki/hls-http.html).

### <i class="fas fa-check"></i> Audio Support
Restreamer recognizes the audio track of your input stream. Optionally you can remove it, add silencs, or convert it to AAC or MP3. [Read more](docs/references-environment-vars.html#rs_audio). 

### <i class="fas fa-check"></i> HTML5 Open Source Video Player
Restreamer uses the [Clappr](https://clappr.io) open source video player. It fully supports HTML5, plays HD, Full-HD, and 4K video, and is responsive.

### <i class="fas fa-check"></i> Snapshots
Restreamer generates JPEG snapshots of your input stream in user defined intervals. [Read more](docs/references-environment-vars.html#rs_snapshot_interval). 

### <i class="fas fa-check"></i> Multiple Streaming
Restreamer works with multiple input streams. Just fire up a new Restreamer instance for each input stream.

### <i class="fas fa-check"></i> Browser Interface
Restreamer lets you easily configure the input stream and the output target via an user interface in your browser.

### <i class="fas fa-check"></i> Password protection
Restreamer protects the browser user interface with a user defined username and password. [Read more](docs/references-environment-vars.html#rs_username).

### <i class="fas fa-check"></i> HTTP API
Restreamer gives you programatically access to the current status of the streams with a JSON HTTP API. [Read more](docs/references-http-api.html).

### <i class="fas fa-check"></i> 100% Data Control
Restreamer guarantees full control of your video data. It doesn't push your video stream anywhere if you don't tell it to. [Read more](docs/guides-external-rtmp.html).

### <i class="fas fa-check"></i> Open Source
Restreamer is open source under the Apache 2.0 license. Check it out on [GitHub](https://github.com/datarhei/restreamer).

### <i class="fas fa-check"></i> Multi Platform Support
Restreamer is built on Docker which allows it to run everywhere, where Docker is supported: Linux, macOS, Windows, Raspberry Pi 1/2/3/Zero, and more. [Read more](docs/installation-index.html).

## Input Video Sources

| Video Source                  |   | Explanation                                       |
|-------------------------------|---|---------------------------------------------------|
| IP/Network Camera             | ✓ | Supporting H.264 - RTP, RTSP, RTMP                |
| Raspicam                      | ✓ | Raspberry Pi Camera Board Module                  |
| USB Camera                    | ✓ | If supported by Linux                             |
| Action Camera                 | ☑ | GoPro had been tested                             |
| Encoder                       | ✓ | Supporting H.264 - RTP, RTSP, RTMP - H.264        |
| H.264                         | ✓ | H.264 sources work                                | 


## Output Video Targets

| Video Destination             |   | Explanation                                       |
|-------------------------------|---|---------------------------------------------------|
| Direct Datarhei streaming     | ✓ | Put your iFrame on your website and start streaming for free|
| YouTube                       | ✓ | Stream live to YouTube.com [howto](../restreamer/docs/guides-youtube.html)|
| Ustream                       | ✓ | Ustream.com                                       |
| Livestream.com                | ✓ | Livestream.com                                    |
| Twitch.tv                     | ✓ | Supporting H.264 - RTP, RTSP, RTMP - H.264        |
| WOWZA                         | ✓ | Push to every WOWZA Streaming Engine [howto](../docs/guides-push-to-wowza.html)|
| Adobe Media Server            | ✓ | Push to Adobe Media Server                        |
| Red5                          | ✓ | Push to open source streaming server Red5         |
| Datarhei Partner              | ✓ | Professional streaming <a target= "_blank" href="http://www.video-stream-hosting.com/">http://www.video-stream-hosting.com/</a>|

## Video Player
Datarhei/Restreamer is using <a target= "_blank" href="http://clappr.io">Clappr.io</a> as video player. 

| Feature                       |   | Explanation                                                 |
|-------------------------------|---|-------------------------------------------------------------|
| HTML5                         | ✓ |                                                             |
| High Resolution 4k Video      | ✓ | All resolutions are supported even 4k and higher            |
| Responsive                    | ✓ |                                                             |
| Audio                         | ✓ |                                                             |
| Watermark                     | ✓ | Put you own corporate Logo in the Player                    |
| Branding                      | ✓ | Cutomize control bar to your company colours                |
| Player Poster                 | ✓ | Snapshot as background if video is stoped or during loading |
| Player size                   | ✓ | Set with width and height paramters in HTML                 |
| Analytics                     | ☑ | Integrate your Google Analytics                             |

## Supported Browsers

| HLS Support                   |   | 
|-------------------------------|---|
| Microsoft Internet Exlorer 10 | ✓ | 
| Microsoft Internet Exlorer 11 | ✓ | 
| Microsoft Edge                | ✓ |  
| Mozilla Firefox               | ✓ |  
| Google Chrome                 | ✓ |  
| Apple Safari                  | ✓ |  
| iPhone                        | ✓ |  
| iPad                          | ✓ |  
| Android                       | ✓ |  
| WiiU Browser                  | ✓ |  
| PS4 Browser                   | ✓ |  


## Support

| Video Player                  |   | Explanation                                                 |
|-------------------------------|---|-------------------------------------------------------------|
| E-Mail                        | ✓ | Ask questions via email (Mo-Fr 24h response time)           |
| Forum on Google Groups        | ✓ | Ask quesions in our support forum 12-24 (Mo-Fr 24h response time)|
| GitHub                        | ✓ | Developer support on GitHub                                 |
| Wiki                          | ✓ | Hundreds of informations round video streaming microcosm    |
