---
title: video-compression
---
######[Camera](../wiki/camera-technology.html) > Video Compression
#Video Compression
## Example to understand video compression
Imagine a green meadow against a cloudless blue sky. On this meadow there is  a cow thatmoves. If the cow is moving, quite simply transfer the image of the cow often. Meadow and sky do not change so often in the frames and so data can be saved. Very simple principle.  

This means that a lot of moving objects or even large areas of water in the picture often require high bandwidth. Even at night, it can  use more bandwidth, but this usually is not too bad, as the camera does hardly show anything due to the lack of light and, luckily, most viewers are not awake at night!  

## What can cause more bandwidth?  
* Large water surfaces, sea, lakes, rivers  
* Trees with many leaves   
* Large crowds  
* Busy streets  
* Night pictures  

--- 
**Datarhei Hint** ☺ ► This are empirical values from 10 years live streaming, it can always happen anyway that a stream much more or less uplaod needed (usually more!).  

---
##Scenarios  
* Rooftops / Still Life / No Movment - less bandwidth  
* Scenics - less bandwidth  
* Street - more bandwidth  
* Beach / sea - a lot of bandwidth  
* Night - a lot of bandwidth  
* River / Water - a lot of bandwidth  

##Rule of thumb for bandwidth in these scenarios (all upload speed) for a static network camera 720p resolution
Rooftops 720p → 100 Kbit/s to 1 Mbit/s  
Scenics 720p → 0.5 Kbit/s - 1.2 Mbit/s  
Street 720p → from 0.5 Kbit/s - 2 Mbit/s  
Beach / Sea → from 2-3 Mbit/s  
Night from 2-3 Mbit/s  
River → from 1.5 - 2 Mbit/s  
Event → 0.5 kbit/s -3 Mbit/s  

**At 450 p or 420p you have about 1/2 less data volume - expected 2-4 times more data volume at 1080p.  **

---
**Datarhei Hint** ☺ ► As a rule of thumb, have at least 1 Mbit/s upload available for a static network camera with 720p stream. 1080p currently plays no major role in live streaming since the user equipment (even on mobile devices via mobile network) often can’t cope with the high bandwidths. For PTZ cameras, we recommend that you have at least 1-2 Mbit/s upload available.  

---
For the more interested of you information on 
<a href="https://en.wikipedia.org/wiki/Video_coding_format" target="_blank">Wikipedia</a>