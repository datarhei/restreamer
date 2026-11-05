package skills

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInvalidBinary(t *testing.T) {
	skills, err := New("")

	require.NotNil(t, err)

	require.Empty(t, skills.FFmpeg.Version)
	require.Empty(t, skills.FFmpeg.Compiler)
	require.Empty(t, skills.FFmpeg.Configuration)
	require.Empty(t, skills.FFmpeg.Libraries)

	require.Empty(t, skills.Filters)
	require.Empty(t, skills.HWAccels)

	require.Empty(t, skills.Codecs.Audio)
	require.Empty(t, skills.Codecs.Subtitle)
	require.Empty(t, skills.Codecs.Video)

	require.Empty(t, skills.Devices.Demuxers)
	require.Empty(t, skills.Devices.Muxers)

	require.Empty(t, skills.Formats.Demuxers)
	require.Empty(t, skills.Formats.Muxers)

	require.Empty(t, skills.Protocols.Input)
	require.Empty(t, skills.Protocols.Output)
}

func TestPatchVersion(t *testing.T) {
	data := `ffmpeg version 4.3.1 Copyright (c) 2000-2020 the FFmpeg developers
 built with Apple clang version 12.0.0 (clang-1200.0.32.29)
 configuration: --enable-static --enable-debug --disable-doc --enable-libx264 --enable-gpl --enable-nonfree
 libavutil      56. 51.100 / 56. 51.100
 libavcodec     58. 91.100 / 58. 91.100
 libavformat    58. 45.100 / 58. 45.100
 libavdevice    58. 10.100 / 58. 10.100
 libavfilter     7. 85.100 /  7. 85.100
 libswscale      5.  7.100 /  5.  7.100
 libswresample   3.  7.100 /  3.  7.100
 libpostproc    55.  7.100 / 55.  7.100`

	f := parseVersion([]byte(data))

	require.Equal(t, ffmpeg{
		Version:       "4.3.1",
		Compiler:      "Apple clang version 12.0.0 (clang-1200.0.32.29)",
		Configuration: "--enable-static --enable-debug --disable-doc --enable-libx264 --enable-gpl --enable-nonfree",
		Libraries: []Library{
			{
				Name:     "libavutil",
				Compiled: "56. 51.100",
				Linked:   "56. 51.100",
			},
			{
				Name:     "libavcodec",
				Compiled: "58. 91.100",
				Linked:   "58. 91.100",
			},
			{
				Name:     "libavformat",
				Compiled: "58. 45.100",
				Linked:   "58. 45.100",
			},
			{
				Name:     "libavdevice",
				Compiled: "58. 10.100",
				Linked:   "58. 10.100",
			},
			{
				Name:     "libavfilter",
				Compiled: "7. 85.100",
				Linked:   "7. 85.100",
			},
			{
				Name:     "libswscale",
				Compiled: "5.  7.100",
				Linked:   "5.  7.100",
			},
			{
				Name:     "libswresample",
				Compiled: "3.  7.100",
				Linked:   "3.  7.100",
			},
			{
				Name:     "libpostproc",
				Compiled: "55.  7.100",
				Linked:   "55.  7.100",
			},
		},
	}, f)
}

func TestMinorVersion(t *testing.T) {
	data := `ffmpeg version 4.4 Copyright (c) 2000-2020 the FFmpeg developers
 built with Apple clang version 12.0.0 (clang-1200.0.32.29)
 configuration: --enable-static --enable-debug --disable-doc --enable-libx264 --enable-gpl --enable-nonfree
 libavutil      56. 51.100 / 56. 51.100
 libavcodec     58. 91.100 / 58. 91.100
 libavformat    58. 45.100 / 58. 45.100
 libavdevice    58. 10.100 / 58. 10.100
 libavfilter     7. 85.100 /  7. 85.100
 libswscale      5.  7.100 /  5.  7.100
 libswresample   3.  7.100 /  3.  7.100
 libpostproc    55.  7.100 / 55.  7.100`

	f := parseVersion([]byte(data))

	require.Equal(t, ffmpeg{
		Version:       "4.4.0",
		Compiler:      "Apple clang version 12.0.0 (clang-1200.0.32.29)",
		Configuration: "--enable-static --enable-debug --disable-doc --enable-libx264 --enable-gpl --enable-nonfree",
		Libraries: []Library{
			{
				Name:     "libavutil",
				Compiled: "56. 51.100",
				Linked:   "56. 51.100",
			},
			{
				Name:     "libavcodec",
				Compiled: "58. 91.100",
				Linked:   "58. 91.100",
			},
			{
				Name:     "libavformat",
				Compiled: "58. 45.100",
				Linked:   "58. 45.100",
			},
			{
				Name:     "libavdevice",
				Compiled: "58. 10.100",
				Linked:   "58. 10.100",
			},
			{
				Name:     "libavfilter",
				Compiled: "7. 85.100",
				Linked:   "7. 85.100",
			},
			{
				Name:     "libswscale",
				Compiled: "5.  7.100",
				Linked:   "5.  7.100",
			},
			{
				Name:     "libswresample",
				Compiled: "3.  7.100",
				Linked:   "3.  7.100",
			},
			{
				Name:     "libpostproc",
				Compiled: "55.  7.100",
				Linked:   "55.  7.100",
			},
		},
	}, f)
}

func TestCustomVersion(t *testing.T) {
	data := `ffmpeg version 4.4.1-datarhei Copyright (c) 2000-2021 the FFmpeg developers
built with gcc 10.3.1 (Alpine 10.3.1_git20211027) 20211027
configuration: --extra-version=datarhei --prefix=/usr --extra-libs='-lpthread -lm -lz -lsupc++ -lstdc++ -lssl -lcrypto -lz -lc -ldl' --enable-nonfree --enable-gpl --enable-version3 --enable-postproc --enable-static --enable-openssl --enable-omx --enable-omx-rpi --enable-mmal --enable-v4l2_m2m --enable-libfreetype --enable-libsrt --enable-libx264 --enable-libx265 --enable-libvpx --enable-libmp3lame --enable-libopus --enable-libvorbis --disable-ffplay --disable-debug --disable-doc --disable-shared
libavutil      56. 70.100 / 56. 70.100
libavcodec     58.134.100 / 58.134.100
libavformat    58. 76.100 / 58. 76.100
libavdevice    58. 13.100 / 58. 13.100
libavfilter     7.110.100 /  7.110.100
libswscale      5.  9.100 /  5.  9.100
libswresample   3.  9.100 /  3.  9.100
libpostproc    55.  9.100 / 55.  9.100`

	f := parseVersion([]byte(data))

	require.Equal(t, ffmpeg{
		Version:       "4.4.1",
		Compiler:      "gcc 10.3.1 (Alpine 10.3.1_git20211027) 20211027",
		Configuration: "--extra-version=datarhei --prefix=/usr --extra-libs='-lpthread -lm -lz -lsupc++ -lstdc++ -lssl -lcrypto -lz -lc -ldl' --enable-nonfree --enable-gpl --enable-version3 --enable-postproc --enable-static --enable-openssl --enable-omx --enable-omx-rpi --enable-mmal --enable-v4l2_m2m --enable-libfreetype --enable-libsrt --enable-libx264 --enable-libx265 --enable-libvpx --enable-libmp3lame --enable-libopus --enable-libvorbis --disable-ffplay --disable-debug --disable-doc --disable-shared",
		Libraries: []Library{
			{
				Name:     "libavutil",
				Compiled: "56. 70.100",
				Linked:   "56. 70.100",
			},
			{
				Name:     "libavcodec",
				Compiled: "58.134.100",
				Linked:   "58.134.100",
			},
			{
				Name:     "libavformat",
				Compiled: "58. 76.100",
				Linked:   "58. 76.100",
			},
			{
				Name:     "libavdevice",
				Compiled: "58. 13.100",
				Linked:   "58. 13.100",
			},
			{
				Name:     "libavfilter",
				Compiled: "7.110.100",
				Linked:   "7.110.100",
			},
			{
				Name:     "libswscale",
				Compiled: "5.  9.100",
				Linked:   "5.  9.100",
			},
			{
				Name:     "libswresample",
				Compiled: "3.  9.100",
				Linked:   "3.  9.100",
			},
			{
				Name:     "libpostproc",
				Compiled: "55.  9.100",
				Linked:   "55.  9.100",
			},
		},
	}, f)
}

func TestFilters(t *testing.T) {
	data := ` ... afirsrc           |->A       Generate a FIR coefficients audio stream.
 ... anoisesrc         |->A       Generate a noise audio signal.
 ... anullsrc          |->A       Null audio source, return empty audio frames.
 ... hilbert           |->A       Generate a Hilbert transform FIR coefficients.
 ... sinc              |->A       Generate a sinc kaiser-windowed low-pass, high-pass, band-pass, or band-reject FIR coefficients.
 ... sine              |->A       Generate sine wave audio signal.
 ... anullsink         A->|       Do absolutely nothing with the input audio.
 ... addroi            V->V       Add region of interest to frame.
 ... alphaextract      V->N       Extract an alpha channel as a grayscale image component.
 T.. alphamerge        VV->V      Copy the luma value of the second input into the alpha channel of the first input.`

	f := parseFilters([]byte(data))

	require.Equal(t, []Filter{
		{
			Id:   "afirsrc",
			Name: "Generate a FIR coefficients audio stream.",
		},
		{
			Id:   "anoisesrc",
			Name: "Generate a noise audio signal.",
		},
		{
			Id:   "anullsrc",
			Name: "Null audio source, return empty audio frames.",
		},
		{
			Id:   "hilbert",
			Name: "Generate a Hilbert transform FIR coefficients.",
		},
		{
			Id:   "sinc",
			Name: "Generate a sinc kaiser-windowed low-pass, high-pass, band-pass, or band-reject FIR coefficients.",
		},
		{
			Id:   "sine",
			Name: "Generate sine wave audio signal.",
		},
		{
			Id:   "anullsink",
			Name: "Do absolutely nothing with the input audio.",
		},
		{
			Id:   "addroi",
			Name: "Add region of interest to frame.",
		},
		{
			Id:   "alphaextract",
			Name: "Extract an alpha channel as a grayscale image component.",
		},
		{
			Id:   "alphamerge",
			Name: "Copy the luma value of the second input into the alpha channel of the first input.",
		},
	}, f)
}

func TestCodecs(t *testing.T) {
	data := ` DEAIL. aac                  AAC (Advanced Audio Coding) (decoders: aac aac_fixed aac_at ) (encoders: aac aac_at )
	DEVI.S y41p                 Uncompressed YUV 4:1:1 12-bit
	DEV.LS h264                 H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10 (encoders: libx264 libx264rgb h264_videotoolbox )
	DEV.L. flv1                 FLV / Sorenson Spark / Sorenson H.263 (Flash Video) (decoders: flv ) (encoders: flv )`

	c := parseCodecs([]byte(data))

	require.Equal(t, ffCodecs{
		Audio: []Codec{
			{
				Id:   "aac",
				Name: "AAC (Advanced Audio Coding)",
				Encoders: []string{
					"aac",
					"aac_at",
				},
				Decoders: []string{
					"aac",
					"aac_fixed",
					"aac_at",
				},
			},
		},
		Video: []Codec{
			{
				Id:   "y41p",
				Name: "Uncompressed YUV 4:1:1 12-bit",
				Encoders: []string{
					"y41p",
				},
				Decoders: []string{
					"y41p",
				},
			},
			{
				Id:   "h264",
				Name: "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
				Encoders: []string{
					"libx264",
					"libx264rgb",
					"h264_videotoolbox",
				},
				Decoders: []string{
					"h264",
				},
			},
			{
				Id:   "flv1",
				Name: "FLV / Sorenson Spark / Sorenson H.263 (Flash Video)",
				Encoders: []string{
					"flv",
				},
				Decoders: []string{
					"flv",
				},
			},
		},
		Subtitle: nil,
	}, c)
}

func TestFormats(t *testing.T) {
	data := ` DE mpeg            MPEG-1 Systems / MPEG program stream
  E mpeg1video      raw MPEG-1 video
  E mpeg2video      raw MPEG-2 video
 DE mpegts          MPEG-TS (MPEG-2 Transport Stream)
 D  mpegtsraw       raw MPEG-TS (MPEG-2 Transport Stream)
 D  mpegvideo       raw MPEG video`

	f := parseFormats([]byte(data))

	require.Equal(t, ffFormats{
		Demuxers: []Format{
			{
				Id:   "mpeg",
				Name: "MPEG-1 Systems / MPEG program stream",
			},
			{
				Id:   "mpegts",
				Name: "MPEG-TS (MPEG-2 Transport Stream)",
			},
			{
				Id:   "mpegtsraw",
				Name: "raw MPEG-TS (MPEG-2 Transport Stream)",
			},
			{
				Id:   "mpegvideo",
				Name: "raw MPEG video",
			},
		},
		Muxers: []Format{
			{
				Id:   "mpeg",
				Name: "MPEG-1 Systems / MPEG program stream",
			},
			{
				Id:   "mpeg1video",
				Name: "raw MPEG-1 video",
			},
			{
				Id:   "mpeg2video",
				Name: "raw MPEG-2 video",
			},
			{
				Id:   "mpegts",
				Name: "MPEG-TS (MPEG-2 Transport Stream)",
			},
		},
	}, f)
}

func TestProtocols(t *testing.T) {
	data := `Input:
	async
	bluray
	cache
Output:
	crypto
	file
	ftp
	gopher`

	p := parseProtocols([]byte(data))

	require.Equal(t, ffProtocols{
		Input: []Protocol{
			{
				Id:   "async",
				Name: "async",
			},
			{
				Id:   "bluray",
				Name: "bluray",
			},
			{
				Id:   "cache",
				Name: "cache",
			},
		},
		Output: []Protocol{
			{
				Id:   "crypto",
				Name: "crypto",
			},
			{
				Id:   "file",
				Name: "file",
			},
			{
				Id:   "ftp",
				Name: "ftp",
			},
			{
				Id:   "gopher",
				Name: "gopher",
			},
		},
	}, p)
}

func TestHWAccels(t *testing.T) {
	data := `Hardware acceleration methods:
videotoolbox`

	p := parseHWAccels([]byte(data))

	require.Equal(t, []HWAccel{
		{
			Id:   "videotoolbox",
			Name: "videotoolbox",
		},
	}, p)
}
