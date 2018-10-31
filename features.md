---
layout: splash
title: Features
header:
  overlay_color: "#5e616c"
  overlay_image: "img/splash.jpg"
excerpt: Restreamer is free open source live video streaming on your website without a streaming provider.
features:
  - title: "H.264 HLS Streaming"
    excerpt: "Supports modern [H.264 streaming](wiki/hls-http.html) from the input than runs natively in your browser without requiring Flash or any other plugins."
  - title: "HTML5 Video Player"
    excerpt: "Open source HTML5 responsive video player that plays HD, Full-HD, and 4K video with audio in fullscreen."
  - title: "Multi Platform"
    excerpt: "Built on Docker it runs on many platforms like Windows, macOS, Linux, Raspberry Pi 1/2/3/Zero, and more."
  - title: "Open Source"
    excerpt: "Restreamer is free and open source (licensed under Apache 2.0), so you can use it for any purpose, private or commercial."
  - title: "Stream Anything"
    excerpt: "Use the H.264 stream from any source such as an [IP cameras](docs/guides-ipcam-rtsp.html), [USB cameras](docs/guides-usb-camera.html), [RaspiCam](docs/guides-raspicam.html), or any H.264 encoder."
  - title: "Stream Anywhere"
    excerpt: "Easily stream your camera to your [website](docs/guides-embed.html), [YouTube](docs/guides-youtube.html), Facebook, twitch.tv, Periscope, Vimeo, [Wowza](docs/guides-wowza.html), AMS, Red5, and more"
  - title: Audio Support
    excerpt: Use the audio track of your input stream and optinally remove, add silence, or convert it to AAC or MP3. [Read more](docs/references-environment-vars.html#rs_audio).
  - title: Snapshots
    excerpt: Generate JPEG snapshots of your input stream in user defined intervals. [Read more](docs/references-environment-vars.html#rs_snapshot_interval). 
  - title: Multiple Streaming
    excerpt: Works with multiple input streams. Just fire up a new Restreamer instance for each input stream.
  - title: Browser Interface
    excerpt: Easily configure the input stream and the output target via an user interface in your browser.
  - title: Password Protection
    excerpt: Password protected browser user interface with a user defined username and password. [Read more](docs/references-environment-vars.html#rs_username).
  - title: HTTP API
    excerpt: Programatically access the current status of the streams with a JSON HTTP API. [Read more](docs/references-http-api.html).
  - title: 100% Data Control
    excerpt: Full control of your video data. The video stream stays with you. [Read more](docs/guides-external-rtmp.html).
  - title: Browser Support
    excerpt: Supports all major browsers on the desktop (IE, Edge, Chrome, Firefox, Safari) as well as on mobile platforms (iPhone, iPad, Android) and game consoles (WiiU, PS4).
  - title: User Support
    excerpt: Support via email, [Google Groups](https://groups.google.com/forum/#!forum/datarhei), [GitHub](https://github.com/datarhei/restreamer/issues), and the [wiki](wiki/index.html).
---

{% include feature_row id="features" %}

### <i class="fas fa-check"></i> Live Video Streaming
Restreamer is an out-of-the-box free video streaming server and ready to stream directly on your website. [Read more](docs/guides-embedding.html).

### <i class="fas fa-check"></i> H.264
Restreamer supports end-to-end H.264 streaming from the input to the target without quality loss.

### <i class="fas fa-check"></i> Stream Anything
Restreamer can use many H.264 input sources for the video such as [IP cameras](docs/guides-ipcam-rtsp.html), [RaspiCam](docs/guides-raspicam.html), [USB cameras](docs/guides-usb-camera.html), action cameras (tested with GoPro), and encoders.

### <i class="fas fa-check"></i> Stream Anywhere
Restreamer makes it easy to push your input video to your [website](docs/guides-embed.html), [YouTube](docs/guides-youtube.html), Facebook, twitch.tv, Periscope, Vimeo, [Wowza](docs/guides-wowza.html), Adobe Media Server, Red5, and more.

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

### <i class="fas fa-check"></i> Password Protection
Restreamer protects the browser user interface with a user defined username and password. [Read more](docs/references-environment-vars.html#rs_username).

### <i class="fas fa-check"></i> HTTP API
Restreamer gives you programatically access to the current status of the streams with a JSON HTTP API. [Read more](docs/references-http-api.html).

### <i class="fas fa-check"></i> 100% Data Control
Restreamer guarantees full control of your video data. It doesn't push your video stream anywhere if you don't tell it to. [Read more](docs/guides-external-rtmp.html).

### <i class="fas fa-check"></i> Open Source
Restreamer is open source under the Apache 2.0 license. Check it out on [GitHub](https://github.com/datarhei/restreamer).

### <i class="fas fa-check"></i> Multi Platform Support
Restreamer is built on Docker which allows it to run everywhere, where Docker is supported: Linux, macOS, Windows, Raspberry Pi 1/2/3/Zero, and more. [Read more](docs/installation-index.html).

### <i class="fas fa-check"></i> Browser Support
Restreamer supports all major browsers on the desktop (IE, Edge, Chrome, Firefox, Safari) as well as on mobile platforms (iPhone, iPad, Android) and game consoles (WiiU, PS4).

### <i class="fas fa-check"></i> User Support
Restreamer offsers support via email, [Google Groups](https://groups.google.com/forum/#!forum/datarhei), [GitHub](https://github.com/datarhei/restreamer/issues), and the [wiki](wiki/index.html).
