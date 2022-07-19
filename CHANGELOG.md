# Restreamer

### 2.2.0

#### Restreamer UI v1.1.0 > v1.2.0

-   Add allow writing HLS to disk
-   Add audio pan filter
-   Add video rotation filter ([#347](https://github.com/datarhei/restreamer/discussions/347))
-   Add video h/v flip filter
-   Add audio volume filter ([#313](https://github.com/datarhei/restreamer/issues/313))
-   Add audio loudness normalization filter
-   Add audio resample filter, that was before part of the encoders
-   Add HLS Master playlist (requires FFmpeg hlsbitrate.patch) (thx Dwaynarang, Electra Player compatibility)
-   Add linkedIn & Azure Media Services to publication services (thx kalashnikov)
-   Add AirPlay support with silvermine videojs plugin
-   Add Chromecast support (thx badincite, [#10](https://github.com/datarhei/restreamer-ui/pull/10))
-   Add stream distribution across multiple internal servers
-   Add SRT settings
-   Add HLS version selection (thx Dwaynarang, Electra Player compatibility)
-   Add Owncast to publication services ([#369](https://github.com/datarhei/restreamer/issues/369))
-   Add Telegram to publication services (thx Martin Held)
-   Add Polish translations (thx Robert Rykała)
-   Mod extends the datarhei Core publication service with srt streaming
-   Mod allow decoders and encoders to set global options
-   Mod allow trailing slash on Core address
-   Fix player problem with different stream formats (9:16)
-   Fix process report naming
-   Fix publication service icon styles
-   Fix VAAPI encoder

#### Core v16.8.0 > v16.9.0

-   Add new placeholders and parameters for placeholder
-   Add optional escape character to process placeholder
-   Add experimental SRT connection stats and logs API
-   Add experimental SRT server (datarhei/gosrt)
-   Add trailing slash for routed directories (datarhei/restreamer#340)
-   Mod allow RTMP server if RTMPS server is enabled
-   Mod create v16 in go.mod
-   Mod allow relative URLs in content in static routes
-   Fix output address validation for tee outputs
-   Fix updating process config
-   Fix hide /config/reload endpoint in reade-only mode
-   Fix data races, tests, lint, and update dependencies

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
