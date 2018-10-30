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
excerpt: 'Live video streaming on your website without a streaming provider.<br /> <small><a href="https://github.com/datarhei/restreamer/releases/tag/0.1.0">Latest release v0.1.0</a></small><br /><br /> {::nomarkdown}<iframe style="display: inline-block;" src="https://ghbtns.com/github-btn.html?user=datarhei&repo=restreamer&type=star&count=true&size=large" frameborder="0" scrolling="0" width="160px" height="30px"></iframe> <iframe style="display: inline-block;" src="https://ghbtns.com/github-btn.html?user=datarhei&repo=restreamer&type=fork&count=true&size=large" frameborder="0" scrolling="0" width="158px" height="30px"></iframe>{:/nomarkdown}'
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
    excerpt: "Supports modern H.264 streaming from the input than runs natively in your browser without requiring Flash or any other plugins."
    url: "features.html"
    btn_class: "btn--primary"
    btn_label: "<i class='far fa-lightbulb'></i> Learn More"
  - title: "HTML5 Video Player"
    excerpt: "Open source HTML5 responsive video player that plays HD, Full-HD, and 4K video with audio in fullscreen."
    url: "features.html"
    btn_class: "btn--primary"
    btn_label: "<i class='far fa-lightbulb'></i> Learn More"
  - title: "Multi Platform"
    excerpt: "Built on Docker it runs on many platforms like Windows, macOS, Linux, Raspberry Pi 1/2/3/Zero, and more."
    url: "features.html"
    btn_class: "btn--primary"
    btn_label: "<i class='far fa-lightbulb'></i> Learn More"
  - title: "Open Source"
    excerpt: "Restreamer is free and open source (licensed under Apache 2.0), so you can use it for any purpose, private or commercial."
    url: "features.html"
    btn_class: "btn--primary"
    btn_label: "<i class='far fa-lightbulb'></i> Learn More"
  - title: "Network Cameras"
    excerpt: "Use the H.264 stream from any source such as an IP camera, USB camera, Rapsicam, or any H.264 encoder."
    url: "features.html"
    btn_class: "btn--primary"
    btn_label: "<i class='far fa-lightbulb'></i> Learn More"
  - title: "Stream Anywhere"
    excerpt: "Easily stream your camera to your website or YouTube, Vimeo, Twitch, Wowza, AMS, Red5, and many more."
    url: "features.html"
    btn_class: "btn--primary"
    btn_label: "<i class='far fa-lightbulb'></i> Learn More"
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
