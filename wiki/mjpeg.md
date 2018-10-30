---
title: MJPEG
---

In Motion JPEG or MJPG full image information is transmitted as a JPG graphic. Unlike H.264 this results in increased data rates for smooth frame
rates. Motion JPG streams are rapidly noticed as "bumpy" for the audience. MJPEG is pure 90ies!
<a href="https://en.wikipedia.org/wiki/Motion_JPEG" target="_blank">Wikipedia</a>

Motion JPEG is supported by Restreamer. But the encoding requires processing power. Raspberry Pi is not suitable for this. Restreamer
relies on [H.264](../wiki/h264.html) via [RTSP](rtsp.html) to allow even at low bandwidths liquid streaming.
{: .notice--info}
