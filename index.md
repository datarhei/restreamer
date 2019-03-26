---
layout: splash
title: Restreamer
permalink: /
header:
  overlay_color: "#5e616c"
  overlay_image: "img/splash.jpg"
  actions:
    - label: "<i class='fab fa-docker'></i> Get Started"
      url: "docs/installation-index.html"
  caption: 
excerpt: 'Live video streaming on your website without a streaming provider.<br /> <small><a href="https://github.com/datarhei/restreamer/releases/tag/v0.2.0">Latest release v0.2.0</a></small><br /><br /> {::nomarkdown}<iframe style="display: inline-block;" src="https://ghbtns.com/github-btn.html?user=datarhei&repo=restreamer&type=star&count=true&size=large" frameborder="0" scrolling="0" width="160px" height="30px"></iframe> <iframe style="display: inline-block;" src="https://ghbtns.com/github-btn.html?user=datarhei&repo=restreamer&type=fork&count=true&size=large" frameborder="0" scrolling="0" width="158px" height="30px"></iframe>{:/nomarkdown}'
intro:
  - title: "Free Video Streaming"
    excerpt: "Restreamer offers smart free video streaming in real time. Stream H.264 video of IP cameras live to your website."
    url: "docs/installation-index.html"
    btn_class: "btn--success"
    btn_label: "<i class='fab fa-docker'></i> Get Started"
  - title: "Distribute"
    excerpt: "Upload your live stream to YouTube, Ustream, Twitch, Livestream.com or any other streaming solutions like Wowza."
    url: "docs/installation-index.html"
    btn_class: "btn--success"
    btn_label: "<i class='fab fa-docker'></i> Get Started"
  - title: "Runs Everywhere"
    excerpt: "The Restreamer Docker image is easy to install and runs on Linux, macOS and Windows, as well as on Raspberry Pi and others."
    url: "docs/installation-index.html"
    btn_class: "btn--success"
    btn_label: "<i class='fab fa-docker'></i> Get Started"
features:
  - title: "H.264 HLS Streaming"
    excerpt: "Supports modern [H.264 streaming](wiki/hls-http.html) from the input than runs natively in your browser without requiring Flash or any other plugins."
  - title: "HTML5 Video Player"
    excerpt: "Open source HTML5 responsive video player that plays HD, Full-HD, and 4K video with audio in fullscreen."
  - title: "Multi Platform"
    excerpt: "Built on Docker it runs on many platforms like [Windows](docs/installation-osx-windows.html), [macOS](docs/installation-osx-windows.html), [Linux](docs/installation-linux-64.html), [Raspberry Pi](docs/installation-linux-arm.html) 1/2/3/Zero, and more."
  - title: "Open Source"
    excerpt: "Restreamer is free and open source (licensed under Apache 2.0), so you can use it for any purpose, private or commercial."
  - title: "Stream Anything"
    excerpt: "Use the H.264 stream from any source such as an [IP cameras](docs/guides-ipcam-rtsp.html), [USB cameras](docs/guides-usb-camera.html), [RaspiCam](docs/guides-raspicam.html), or any H.264 encoder."
  - title: "Stream Anywhere"
    excerpt: "Easily stream your camera to your [website](docs/guides-embedding.html), [YouTube](docs/guides-youtube.html), Facebook, twitch.tv, Periscope, Vimeo, [Wowza](docs/guides-wowza.html), AMS, Red5, and more"
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
install:
  - title: "Windows / macOS"
    excerpt: "Install on Windows or macOS with Docker and Kitematic in just a few clicks."
    url: "docs/installation-osx-windows.html"
    btn_class: "btn--primary"
    btn_label: "<i class='fa fa-cloud-download-alt'></i> Install"
  - title: "Linux"
    excerpt: "Install on Linux with a 64bit Intel or AMD CPU using Docker."
    url: "docs/installation-linux-64.html"
    btn_class: "btn--primary"
    btn_label: "<i class='fa fa-cloud-download-alt'></i> Install"
  - title: "Raspberry Pi / ARM"
    excerpt: "Install on Raspberry Pi 1, Pi 2, and Pi 3 or any other device with an ARMv6, ARMv7, or ARMv8 CPU."
    url: "docs/installation-linux-arm.html"
    btn_class: "btn--primary"
    btn_label: "<i class='fa fa-cloud-download-alt'></i> Install"
powered:
  - title: "Powered by"
    excerpt: "![Intro Powered By](img/intro_powered.png)"
---

{% include feature_row id="intro" %}

![Intro Banner](img/intro_banner.png){: .align-center}

{% include feature_row id="features" %}

![Intro Docker](img/intro_docker.png){: .align-center}

{% include feature_row id="install" %}

{% include feature_row id="powered" type="center" %}
