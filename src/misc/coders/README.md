# Decoders and Encoders

Implementations of various decoders and encoder for audio and video.

## Decoders

Each decoder exports the same variables:

```
export { coder, name, codecs, type, hwaccel, defaults, Coder as component };
```

| Variable | Description                                                          |
| -------- | -------------------------------------------------------------------- |
| coder    | Name of the decoder in FFmpeg, e.g. `cuda`, `vc1_mmal`.              |
| name     | Name for the decoder as it will be displayed in the UI.              |
| codecs   | Array of codecs this coder supports, e.g. `['h264', 'h265']`.        |
| type     | Either `video` or `audio`.                                           |
| hwaccel  | Whether this codec uses hardware acceleration.                       |
| defaults | A function that returns the default settings and mapping. See below. |
| Coder    | The React component.                                                 |

### defaults

The `defaults(stream, skills)` function returns the default settings and mappings for this decoder. It is an object of this shape:

```
{
	settings: {},
	mapping: {
        global: [],
        local: [],
    },
}
```

The `settings` is an object private to a coder and contains its settings as required for rendering the component
with options for this coder. The `mapping` object contains the FFmpeg command line options according to the settings.
It has a `global` array which contains all global options for this coder. _Each option (with its value) has to be
an array of its own_. The `local` array is an array of options for that input, e.g.

```
{
	settings: {
        ...
    },
	mapping: {
        global: [
            ['-init_hw_device', 'vaapi=foo:/dev/dri/renderD128'],
        ],
        local: [
            '-hwaccel', 'vaapi',
            '-hwaccel_output_format', 'vaapi',
            '-hwaccel_device', 'foo',
        ],
    },
}
```

Check out the existing decoders as examples for an implementation.

## Encoders

Each encoder exports the same variables:

```
export { coder, name, codec, type, hwaccel, summarize, defaults, Coder as component };
```

| Variable  | Description                                                            |
| --------- | ---------------------------------------------------------------------- |
| coder     | Name of the encoder in FFmpeg, e.g. `libx264`.                         |
| name      | Name for the encoder as it will be displayed in the UI.                |
| codec     | Name of the codec, e.g. `h264`.                                        |
| type      | Either `video` or `audio`.                                             |
| hwaccel   | Whether this codec uses hardware acceleration.                         |
| summarize | A function that returns a string that summarizes the current settings. |
| defaults  | A function that returns the default settings and mapping. See below.   |
| Coder     | The React component.                                                   |

### defaults

The `defaults(stream, skills)` function returns the default settings and mappings for this encoder. It is an object of this shape:

```
{
	settings: {},
	mapping: {
        global: [],
        local: [],
    },
}
```

The `settings` is an object private to a coder and contains its settings as required for rendering the component
with options for this coder. The `mapping` object contains the FFmpeg command line options according to the settings.
It has a `global` array which contains all global options for this coder. _Each option (with its value) has to be
an array of its own_. The `local` array is an array of options for that input, e.g.

```
{
	settings: {
        ...
    },
	mapping: {
        global: [
            ['-init_hw_device', 'vaapi=foo:/dev/dri/renderD128'],
        ],
        local: [
            '-filter_hw_device', 'foo',
            '-filter:v', 'format=nv12|vaapi,hwupload',
            '-codec:v', 'h264_vaapi',
        ],
    },
}
```

Check out the existing encoders as examples for an implementation.
