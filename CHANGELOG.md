## Changes from 0.1.0-RC3 to 0.1.0-RC4

* added https sources to Dockerfiles (thx @ [nodiscc](https://github.com/nodiscc))
* fixed visualisation problem of new RTMP/RTSP endpoint (JSONDB processing)
* optimized input-process (removed generateOutputRTMPPath)
* cleaned up NPM packages
* added ESLint insteed of JSHint+ first ruleset
* first code optimizations (ecma5 -> ecma6)
* optimized RaspiCam-Hack for Raspberry Pi 1
* updated NPM packages and node
* refactored FFmpeg-Fluent integration by own fork (reduces process-output)
* fixed output-process (generated high CPU load on arm)
* fixed preview player (no play-icon, didn't stop when the modalbox is closed)

##### Merged pull requests:

* Fixed ARM typo @ [Vyacheslav Shevchenko](https://github.com/bliz937)

## Changes from 0.1.0-RC2 to 0.1.0-RC3

* fixed links to docs
* updated NPM & Bower dependencies

## Changes from 0.1.0-RC1 to 0.1.0-RC2

* use branch FIX_stderr-variable_of_infinite_growth
* more jshint stylecheck rules and fix problems
* renamed ReStreamer to Restreamer
* refactored docs and wiki
