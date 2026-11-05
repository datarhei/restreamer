# Patches

FFmpeg patches.

## HLS

Calculate bandwidth estimate for the hls master manifest.

## JSON stats

This repository contains a patch for the FFmpeg program to provide detailed progress information. With this patch, FFmpeg will output
the progress information in a JSON string that contains the data for each input and output stream individually. The JSON output is enabled
by default. It can be enabled with the global `-jsonstats` switch on the command line. Use the `-stats` switch
on the command line for the standard progress output.

Example output with `-stats`:

```
frame=  143 fps= 25 q=-1.0 Lsize=     941kB time=00:00:05.68 bitrate=1357.0kbits/s speed=0.995x
```

Example output with `-jsonstats`:

```
JSONProgress:{"inputs":[{"id":0, "stream":0, "type":"video", "codec":"rawvideo", "coder":"rawvideo", "pix_fmt":"rgb24", "frame":188, "fps":24.95, "width":1280, "height":720, "size_kb":507600, "bitrate_kbps":552960.0},{"id":1, "stream":0, "type":"audio", "codec":"pcm_u8", "coder":"pcm_u8", "frame":314, "sampling_hz":44100, "layout":"stereo", "size_kb":628, "bitrate_kbps":705.6}], "outputs":[{"id":0, "stream":0, "type":"video", "codec":"h264", "coder":"libx264", "pix_fmt":"yuv420p", "frame":188, "fps":24.95, "q":-1.0, "width":1280, "height":720, "size_kb":1247, "bitrate_kbps":1365.6},{"id":0, "stream":1, "type":"audio", "codec":"aac", "coder":"aac", "frame":315, "sampling_hz":44100, "layout":"stereo", "size_kb":2, "bitrate_kbps":2.1}], "frame":188, "fps":24.95, "q":-1.0, "size_kb":1249, "bitrate_kbps":1367.7, "time":"0h0m7.48s", "speed":0.993, "dup":0, "drop":0}
```

The same output but nicely formatted:

```json
{
  "bitrate_kbps": 1367.7,
  "drop": 0,
  "dup": 0,
  "fps": 24.95,
  "frame": 188,
  "inputs": [
    {
      "bitrate_kbps": 552960.0,
      "codec": "rawvideo",
      "coder": "rawvideo",
      "fps": 24.95,
      "frame": 188,
      "height": 720,
      "id": 0,
      "pix_fmt": "rgb24",
      "size_kb": 507600,
      "stream": 0,
      "type": "video",
      "width": 1280
    },
    {
      "bitrate_kbps": 705.6,
      "codec": "pcm_u8",
      "coder": "pcm_u8",
      "frame": 314,
      "id": 1,
      "layout": "stereo",
      "sampling_hz": 44100,
      "size_kb": 628,
      "stream": 0,
      "type": "audio"
    }
  ],
  "outputs": [
    {
      "bitrate_kbps": 1365.6,
      "codec": "h264",
      "coder": "libx264",
      "fps": 24.95,
      "frame": 188,
      "height": 720,
      "id": 0,
      "pix_fmt": "yuv420p",
      "q": -1.0,
      "size_kb": 1247,
      "stream": 0,
      "type": "video",
      "width": 1280
    },
    {
      "bitrate_kbps": 2.1,
      "codec": "aac",
      "coder": "aac",
      "frame": 315,
      "id": 0,
      "layout": "stereo",
      "sampling_hz": 44100,
      "size_kb": 2,
      "stream": 1,
      "type": "audio"
    }
  ],
  "q": -1.0,
  "size_kb": 1249,
  "speed": 0.993,
  "time": "0h0m7.48s"
}
```
