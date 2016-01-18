---
title: iframe
---
######[Camera](../wiki/camera-technology.html) > iFrame / GOP / GOV
#iFrame / GOP / GOV
In the field of video compression a video frame is compressed using different algorithms with different advantages and disadvantages, centered mainly around amount of data compression. These different algorithms for video frames are called picture types or frame types. The three major picture types used in the different video algorithms are I, P and B. They are different in the following characteristics:  
I‑frames are the least compressible but don't require other video frames to decode.  
P‑frames can use data from previous frames to decompress and are more compressible than I‑frames.  
B‑frames can use both previous and forward frames for data reference to get the highest amount of data compression. <a href="https://en.wikipedia.org/wiki/Video_compression_picture_types" target="_blank">Wikipedia</a>

**Datarhei Hint ☺** ► In this context also dive often the terms <a href="https://en.wikipedia.org/wiki/Group_of_pictures" target="_blank">GOP (Group of Pictures) or on GOV</a>. The higher the parameter in your ip-camera the lower is the bandwidth required but the lower the quality of your video gets. We recommend an optimal GOP of 25 but no more than 61st.  