# Restreamer

### v2.1.0

-   Fix Dockerfile (bundles frontend, backend and FFmpeg)

#### Restreamer UI v1.0.0 > v1.1.0

-   Add "HLS cleanup" as an optional function ([Philipp Trenz](https://github.com/philipptrenz))
-   Add /ui info to / ([#326](https://github.com/datarhei/restreamer/issues/326))
-   Add Russian translation (thx Inthegamelp)
-   Add missed VAAPI encoder
-   Add missed V4L2_M2M encoder
-   Add missed Raspberry Pi 64bit Docker image
-   Mod updates VideoJS
-   Add option to disable playersites share-button (thx Anders Mellgren)
-   Fix hides unset content license on playersite (thx Anders Mellgren)
-   Fix updates V4L2 device-list on select
-   Fix snapshot interval ([#341](https://github.com/datarhei/restreamer/issues/340))
-   Fix reverse proxy issue ([#340](https://github.com/datarhei/restreamer/issues/340))
-   Fix double escape failer ([#336](https://github.com/datarhei/restreamer/issues/336))
-   Fix type in player plugin ([#336](https://github.com/datarhei/restreamer/issues/336))
-   Fix deletes processes with dependencies (thx Patron Ramakrishna Chillara)
-   Fix datarhei Core publication service
-   Fix dependabot alerts
-   Fix code scanning alerts
-   Merge security pr

Preparation for FFmpeg v5.0 (migration will not work)

-   Add FFmpeg v5.0 commands (preparation)
-   Mod allows FFmpeg v5.0 (preparation)

#### Core v16.7.2 > v16.8.0

-   Add purge_on_delete function
-   Mod updated dependencies
-   Mod updated API docs
-   Fix disabled session logging
-   Fix FFmpeg skills reload
-   Fix ignores processes with invalid references (thx Patron Ramakrishna Chillara)
-   Fix code scanning alerts

#### FFmpeg v4.2.2 > v4.2.2-1

-   Mod enables libv4l2 + v4l2_m2m for all Docker images ([#339](https://github.com/datarhei/restreamer/issues/339))
-   Mod updates packages (freetype, libxml2, VAAPI-Libs (gmmlib, media-driver, media-sdk))
-   Fix cuda Docker image ([#328](https://github.com/datarhei/restreamer/issues/328))
-   Fix VAAPI Docker image

#### FFmpeg v5.0.1

-   Add FFmpeg v5.0 JSONSTATS patch (preparation)

#### Documentation

-   Add [Basic troubleshooting guide](https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/basic-troubleshooting)
-   Add [custom playersite guide](https://docs.datarhei.com/restreamer/knowledge-base/user-guides/how-to-integrate-a-website) ([#337](https://github.com/datarhei/restreamer/issues/337))
-   Add [guide for custom RTMP port](https://docs.datarhei.com/restreamer/knowledge-base/user-guides/how-to-change-the-rtmp-port)
-   Add [encoding compatiblity list](https://docs.datarhei.com/restreamer/knowledge-base/troubleshooting/encoding-compatibility-list))

---
