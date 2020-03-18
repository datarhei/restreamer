---
title: Ingest RTMP
---

Inside the Restreamer Docker container is a RTMP-Server running (nginx-rtmp) that you can use to push streams to. To have access to the
RTMP-Server, you have to expose the port 1935 of the Docker container, e.g.

```
docker run ... -p 1935:1935 ... datarhei/restreamer
```

This is inheritly dangerous because anyone can now push RTMP streams to your box. In order to protect the RTMP-Server you should
define a token ([more details](references-environment-vars.html#rs_token)), e.g.

```
docker run ... -p 1935:1935 -e RS_TOKEN=a-secret-only-you-know ... datarhei/restreamer
```

Now you can push a RTMP-Stream to, e.g. `rtmp://[address of your box]/live/external.stream?token=...`.

Navigate to the Restreamer web interface, e.g. `http://[address of your box]:8080/`. For the video source you enter the same address
you are pushing to, and replace the public address and port with localhost, i.e. `rtmp://localhost/live/external.stream?token=...`.

Now you can embed the player in a website or forward it to an [external video streaming provider](guides-external-rtmp.html) like
[YouTube](giudes-youtube.html), [Twitch](guides-twitch.html), Facebook, or others.
