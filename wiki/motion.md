---
title: Motion
---

Motion, a software motion detector, is a free, open source CCTV software application developed for Linux.

It can monitor video signal from one or more cameras and is able to detect if a significant part of the picture has changed saving away video
when it detects that motion is occurring (it can also do time lapse videos, et al.).

The program is written in C and is made for Linux (exploiting video4linux interface). Motion is a command line based tool whose output can be
either jpeg, netpbm files or mpeg video sequences. It is strictly command line driven and can run as a daemon with a rather small footprint
and low CPU usage.

It is operated mainly via config files, though the end video streams can be viewed from a web browser. It can also call to user configurable
"triggers" when certain events occur. <a href="https://en.wikipedia.org/wiki/Motion_%28surveillance_software%29" target="_blank">Wikipedia</a>