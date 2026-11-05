# Core

### Core v16.15.0 > v16.16.0

-   Add ConnectionIdleTimeout to RTMP server
-   Add WithLevel() to Logger interface
-   Fix datarhei/restreamer#759
-   Fix various RTMP bugs
-   Fix wrong log output when receiving a RTMP stream
-   Fix skipping session handling if collectors are nil
-   Update dependencies

### Core v16.14.0 > v16.15.0

-   Add migrating to ffmpeg 6
-   Fix missing process data if process has been deleted meanwhile
-   Fix maintaining the metadata on process config update (datarhei/restreamer#698)
-   Fix placeholder parsing
-   Fix concurrent memfs accesses
-   Fix memfs concurrent read and write performance

### Core v16.13.1 > v16.14.0

-   Add support for SRTv4 clients
-   Add support for Enhanced RTMP in internal RTMP server
-   Fix require positive persist interval (session)
-   Fix race condition (process)
-   Update dependencies

### Core v16.13.0 > v16.13.1

-   Fix transfer of reports to updated process
-   Fix calling Wait after process has been read
-   Fix 509 return code if non-existing stream is requested
-   Fix default search paths for config file
-   Fix sized filesystem
-   Update dependencies

### Core v16.12.0 > v16.13.0

-   Add updated_at field in process infos
-   Add preserve process log history when updating a process
-   Add support for input framerate data from jsonstats patch
-   Add number of keyframes and extradata size to process progress data
-   Mod bumps FFmpeg to v5.1.3 (datarhei/core:tag bundles)
-   Fix better naming for storage endpoint documentation
-   Fix freeing up S3 mounts
-   Fix URL validation if the path contains FFmpeg specific placeholders
-   Fix purging default file from HTTP cache
-   Fix parsing S3 storage definition from environment variable
-   Fix checking length of CPU time array ([#10](https://github.com/datarhei/core/issues/10))
-   Fix possible infinite loop with HLS session rewriter
-   Fix not propagating process limits
-   Fix URL validation if the path contains FFmpeg specific placeholders
-   Fix RTMP DoS attack (thx Johannes Frank)
-   Deprecate ENV names that do not correspond to JSON name

### Core v16.11.0 > v16.12.0

-   Add S3 storage support
-   Add support for variables in placeholde parameter
-   Add support for RTMP token as stream key as last element in path
-   Add support for soft memory limit with debug.memory_limit_mbytes in config
-   Add support for partial process config updates
-   Add support for alternative syntax for auth0 tenants as environment variable
-   Fix config timestamps created_at and loaded_at
-   Fix /config/reload return type
-   Fix modifying DTS in RTMP packets ([restreamer/#487](https://github.com/datarhei/restreamer/issues/487), [restreamer/#367](https://github.com/datarhei/restreamer/issues/367))
-   Fix default internal SRT latency to 20ms

### Core v16.10.1 > v16.11.0

-   Add FFmpeg 4.4 to FFmpeg 5.1 migration tool
-   Add alternative SRT streamid
-   Mod bump FFmpeg to v5.1.2 (datarhei/core:tag bundles)
-   Fix crash with custom SSL certificates ([restreamer/#425](https://github.com/datarhei/restreamer/issues/425))
-   Fix proper version handling for config
-   Fix widged session data
-   Fix resetting process stats when process stopped
-   Fix stale FFmpeg process detection for streams with only audio
-   Fix wrong return status code ([#6](https://github.com/datarhei/core/issues/6)))
-   Fix use SRT defaults for key material exchange

### Core v16.10.0 > v16.10.1

-   Add email address in TLS config for Let's Encrypt
-   Fix use of Let's Encrypt production CA

### Core v16.9.1 > v16.10.0

-   Add HLS session middleware to diskfs
-   Add /v3/metrics (get) endpoint to list all known metrics
-   Add logging HTTP request and response body sizes
-   Add process id and reference glob pattern matching
-   Add cache block list for extensions not to cache
-   Mod exclude .m3u8 and .mpd files from disk cache by default
-   Mod replaces x/crypto/acme/autocert with caddyserver/certmagic
-   Mod exposes ports (Docker desktop)
-   Fix assigning cleanup rules for diskfs
-   Fix wrong path for swagger definition
-   Fix process cleanup on delete, remove empty directories from disk
-   Fix SRT blocking port on restart (upgrade datarhei/gosrt)
-   Fix RTMP communication (Blackmagic Web Presenter, thx 235 MEDIA)
-   Fix RTMP communication (Blackmagic ATEM Mini, [#385](https://github.com/datarhei/restreamer/issues/385))
-   Fix injecting commit, branch, and build info
-   Fix API metadata endpoints responses

#### Core v16.9.0 > v16.9.1^

-   Fix v1 import app
-   Fix race condition

#### Core v16.8.0 > v16.9.0

-   Add new placeholders and parameters for placeholder
-   Allow RTMP server if RTMPS server is enabled. In case you already had RTMPS enabled it will listen on the same port as before. An RTMP server will be started additionally listening on a lower port number. The RTMP app is required to start with a slash.
-   Add optional escape character to process placeholder
-   Fix output address validation for tee outputs
-   Fix updating process config
-   Add experimental SRT connection stats and logs API
-   Hide /config/reload endpoint in reade-only mode
-   Add experimental SRT server (datarhei/gosrt)
-   Create v16 in go.mod
-   Fix data races, tests, lint, and update dependencies
-   Add trailing slash for routed directories (datarhei/restreamer#340)
-   Allow relative URLs in content in static routes

#### Core v16.7.2 > v16.8.0

-   Add purge_on_delete function
-   Mod updated dependencies
-   Mod updated API docs
-   Fix disabled session logging
-   Fix FFmpeg skills reload
-   Fix ignores processes with invalid references (thx Patron Ramakrishna Chillara)
-   Fix code scanning alerts
