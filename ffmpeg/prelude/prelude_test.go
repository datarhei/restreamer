package prelude

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrelude(t *testing.T) {
	rawdata := `ffmpeg version 4.0.2 Copyright (c) 2000-2018 the FFmpeg developers
  built with Apple LLVM version 9.1.0 (clang-902.0.39.2)
  configuration: --prefix=/usr/local/Cellar/ffmpeg/4.0.2 --enable-shared --enable-pthreads --enable-version3 --enable-hardcoded-tables --enable-avresample --cc=clang --host-cflags= --host-ldflags= --enable-gpl --enable-libmp3lame --enable-libx264 --enable-libx265 --enable-libxvid --enable-opencl --enable-videotoolbox --disable-lzma
  libavutil      56. 14.100 / 56. 14.100
  libavcodec     58. 18.100 / 58. 18.100
  libavformat    58. 12.100 / 58. 12.100
  libavdevice    58.  3.100 / 58.  3.100
  libavfilter     7. 16.100 /  7. 16.100
  libavresample   4.  0.  0 /  4.  0.  0
  libswscale      5.  1.100 /  5.  1.100
  libswresample   3.  1.100 /  3.  1.100
  libpostproc    55.  1.100 / 55.  1.100
Input #0, lavfi, from 'testsrc=size=1280x720:rate=25':
  Duration: N/A, start: 0.000000, bitrate: N/A
    Stream #0:0: Video: rawvideo (RGB[24] / 0x18424752), rgb24, 1280x720 [SAR 1:1 DAR 16:9], 25 tbr, 25 tbn, 25 tbc
Input #1, lavfi, from 'anullsrc=r=44100:cl=stereo':
  Duration: N/A, start: 0.000000, bitrate: 705 kb/s
    Stream #1:0: Audio: pcm_u8, 44100 Hz, stereo, u8, 705 kb/s
Input #2, playout, from 'playout:rtmp://l5gn74l5-vpu.livespotting.com/live/0chl6hu7_360?token=m5ZuiCQYRlIon8':
  Duration: N/A, start: 0.000000, bitrate: 265 kb/s
    Stream #2:0: Video: h264 (Constrained Baseline), yuvj420p(pc, progressive), 640x360 [SAR 1:1 DAR 16:9], 265 kb/s, 10 fps, 10 tbr, 1000k tbn, 20 tbc
Input #3, mov,mp4,m4a,3gp,3g2,mj2, from 'movie.mp4':
  Metadata:
    major_brand     : isom
    minor_version   : 512
    compatible_brands: isomiso2avc1mp41
    encoder         : Lavf58.20.100
  Duration: 00:01:02.28, start: 0.000000, bitrate: 5895 kb/s
    Stream #3:0(eng): Video: h264 (Main) (avc1 / 0x31637661), yuvj420p(pc, bt709), 2560x1440 [SAR 1:1 DAR 16:9], 5894 kb/s, 23.98 fps, 25 tbr, 90k tbn, 50 tbc (default)
    Stream #3:1(por): Subtitle: subrip
Input #4, mpegts, from 'srt://localhost:6000?mode=caller&transtype=live&streamid=#!:m=request,r=ingest/ad045490-8233-4f31-a296-ea5771a340ac&passphrase=foobarfoobar':
  Duration: N/A, start: 71.786667, bitrate: N/A
	Program 1
		Metadata:
			service_name    : Service01
			service_provider: FFmpeg
	Stream #4:0[0x100]: Video: h264 (Main) ([27][0][0][0] / 0x001B), yuv420p(tv, smpte170m/bt709/bt709, progressive), 1920x1080 [SAR 1:1 DAR 16:9], 25 tbr, 90k tbn
	Stream #4:1[0x101]: Audio: aac (LC) ([15][0][0][0] / 0x000F), 48000 Hz, stereo, fltp, 162 kb/s
Stream mapping:
  Stream #0:0 -> #0:0 (rawvideo (native) -> h264 (libx264))
  Stream #1:0 -> #0:1 (pcm_u8 (native) -> aac (native))
Press [q] to stop, [?] for help
[libx264 @ 0x7fa96a800600] using SAR=1/1
[libx264 @ 0x7fa96a800600] using cpu capabilities: MMX2 SSE2Fast SSSE3 SSE4.2 AVX FMA3 BMI2 AVX2
[libx264 @ 0x7fa96a800600] profile Constrained Baseline, level 3.1
[libx264 @ 0x7fa96a800600] 264 - core 152 r2854 e9a5903 - H.264/MPEG-4 AVC codec - Copyleft 2003-2017 - http://www.videolan.org/x264.html - options: cabac=0 ref=1 deblock=0:0:0 analyse=0:0 me=dia subme=0 psy=1 psy_rd=1.00:0.00 mixed_ref=0 me_range=16 chroma_me=1 trellis=0 8x8dct=0 cqm=0 deadzone=21,11 fast_pskip=1 chroma_qp_offset=0 threads=6 lookahead_threads=1 sliced_threads=0 nr=0 decimate=1 interlaced=0 bluray_compat=0 constrained_intra=0 bframes=0 weightp=0 keyint=50 keyint_min=5 scenecut=0 intra_refresh=0 rc=crf mbtree=0 crf=23.0 qcomp=0.60 qpmin=0 qpmax=69 qpstep=4 ip_ratio=1.40 aq=0
[hls @ 0x7fa969803a00] Opening './data/testsrc5375.ts.tmp' for writing
Output #0, hls, to './data/testsrc.m3u8':
  Metadata:
    encoder         : Lavf58.12.100
    Stream #0:0: Video: h264 (libx264), yuv420p(progressive), 1280x720 [SAR 1:1 DAR 16:9], q=-1--1, 25 fps, 90k tbn, 25 tbc
    Metadata:
      encoder         : Lavc58.18.100 libx264
    Side data:
      cpb: bitrate max/min/avg: 0/0/0 buffer size: 0 vbv_delay: -1
    Stream #0:1: Audio: aac (LC), 44100 Hz, stereo, fltp, 64 kb/s
    Metadata:
      encoder         : Lavc58.18.100 aac
[hls @ 0x7fa969803a00] Opening './data/testsrc5376.ts.tmp' for writing=0.872x
[hls @ 0x7fa969803a00] Opening './data/testsrc.m3u8.tmp' for writing
[hls @ 0x7fa969803a00] Opening './data/testsrc.m3u8.tmp' for writing
frame=   58 fps= 25 q=-1.0 Lsize=N/A time=00:00:02.32 bitrate=N/A speed=0.999x`

	data := strings.Split(rawdata, "\n")

	inputs, outputs, noutputs := Parse(data)

	require.Equal(t, 7, len(inputs))
	require.Equal(t, 2, len(outputs))
	require.Equal(t, 2, noutputs)

	i := inputs[0]

	require.Equal(t, "testsrc=size=1280x720:rate=25", i.Address)
	require.Equal(t, "lavfi", i.Format)
	require.Equal(t, uint64(0), i.Index)
	require.Equal(t, uint64(0), i.Stream)
	require.Equal(t, "und", i.Language)
	require.Equal(t, "video", i.Type)
	require.Equal(t, "rawvideo", i.Codec)
	require.Equal(t, 0.0, i.Bitrate)
	require.Equal(t, 0.0, i.Duration)
	require.Equal(t, 0.0, i.FPS)
	require.Equal(t, "rgb24", i.Pixfmt)
	require.Equal(t, uint64(1280), i.Width)
	require.Equal(t, uint64(720), i.Height)

	i = inputs[1]

	require.Equal(t, "anullsrc=r=44100:cl=stereo", i.Address)
	require.Equal(t, "lavfi", i.Format)
	require.Equal(t, uint64(1), i.Index)
	require.Equal(t, uint64(0), i.Stream)
	require.Equal(t, "und", i.Language)
	require.Equal(t, "audio", i.Type)
	require.Equal(t, "pcm_u8", i.Codec)
	require.Equal(t, 705.0, i.Bitrate)
	require.Equal(t, 0.0, i.Duration)
	require.Equal(t, uint64(44100), i.Sampling)
	require.Equal(t, "stereo", i.Layout)

	i = inputs[2]

	require.Equal(t, "playout:rtmp://l5gn74l5-vpu.livespotting.com/live/0chl6hu7_360?token=m5ZuiCQYRlIon8", i.Address)
	require.Equal(t, "playout", i.Format)
	require.Equal(t, uint64(2), i.Index)
	require.Equal(t, uint64(0), i.Stream)
	require.Equal(t, "und", i.Language)
	require.Equal(t, "video", i.Type)
	require.Equal(t, "h264", i.Codec)
	require.Equal(t, 265.0, i.Bitrate)
	require.Equal(t, 0.0, i.Duration)
	require.Equal(t, 10.0, i.FPS)
	require.Equal(t, "yuvj420p", i.Pixfmt)
	require.Equal(t, uint64(640), i.Width)
	require.Equal(t, uint64(360), i.Height)

	i = inputs[3]

	require.Equal(t, "movie.mp4", i.Address)
	require.Equal(t, "mov,mp4,m4a,3gp,3g2,mj2", i.Format)
	require.Equal(t, uint64(3), i.Index)
	require.Equal(t, uint64(0), i.Stream)
	require.Equal(t, "eng", i.Language)
	require.Equal(t, "video", i.Type)
	require.Equal(t, "h264", i.Codec)
	require.Equal(t, 5894.0, i.Bitrate)
	require.Equal(t, 62.28, i.Duration)
	require.Equal(t, 23.98, i.FPS)
	require.Equal(t, "yuvj420p", i.Pixfmt)
	require.Equal(t, uint64(2560), i.Width)
	require.Equal(t, uint64(1440), i.Height)

	i = inputs[4]

	require.Equal(t, "movie.mp4", i.Address)
	require.Equal(t, "mov,mp4,m4a,3gp,3g2,mj2", i.Format)
	require.Equal(t, uint64(3), i.Index)
	require.Equal(t, uint64(1), i.Stream)
	require.Equal(t, "por", i.Language)
	require.Equal(t, "subtitle", i.Type)
	require.Equal(t, "subrip", i.Codec)

	i = inputs[5]

	require.Equal(t, "srt://localhost:6000?mode=caller&transtype=live&streamid=#!:m=request,r=ingest/ad045490-8233-4f31-a296-ea5771a340ac&passphrase=foobarfoobar", i.Address)
	require.Equal(t, "mpegts", i.Format)
	require.Equal(t, uint64(4), i.Index)
	require.Equal(t, uint64(0), i.Stream)
	require.Equal(t, "und", i.Language)
	require.Equal(t, "video", i.Type)
	require.Equal(t, "h264", i.Codec)
	require.Equal(t, 0.0, i.Bitrate)
	require.Equal(t, 0.0, i.Duration)
	require.Equal(t, 0.0, i.FPS)
	require.Equal(t, "yuv420p", i.Pixfmt)
	require.Equal(t, uint64(1920), i.Width)
	require.Equal(t, uint64(1080), i.Height)

	i = inputs[6]

	require.Equal(t, "srt://localhost:6000?mode=caller&transtype=live&streamid=#!:m=request,r=ingest/ad045490-8233-4f31-a296-ea5771a340ac&passphrase=foobarfoobar", i.Address)
	require.Equal(t, "mpegts", i.Format)
	require.Equal(t, uint64(4), i.Index)
	require.Equal(t, uint64(1), i.Stream)
	require.Equal(t, "und", i.Language)
	require.Equal(t, "audio", i.Type)
	require.Equal(t, "aac", i.Codec)
	require.Equal(t, 162.0, i.Bitrate)
	require.Equal(t, 0.0, i.Duration)
	require.Equal(t, uint64(48000), i.Sampling)
	require.Equal(t, "stereo", i.Layout)

	i = outputs[0]

	require.Equal(t, "./data/testsrc.m3u8", i.Address)
	require.Equal(t, "hls", i.Format)
	require.Equal(t, uint64(0), i.Index)
	require.Equal(t, uint64(0), i.Stream)
	require.Equal(t, "und", i.Language)
	require.Equal(t, "video", i.Type)
	require.Equal(t, "h264", i.Codec)
	require.Equal(t, 0.0, i.Bitrate)
	require.Equal(t, 0.0, i.Duration)
	require.Equal(t, 25.0, i.FPS)
	require.Equal(t, "yuv420p", i.Pixfmt)
	require.Equal(t, uint64(1280), i.Width)
	require.Equal(t, uint64(720), i.Height)

	i = outputs[1]

	require.Equal(t, "./data/testsrc.m3u8", i.Address)
	require.Equal(t, "hls", i.Format)
	require.Equal(t, uint64(0), i.Index)
	require.Equal(t, uint64(1), i.Stream)
	require.Equal(t, "und", i.Language)
	require.Equal(t, "audio", i.Type)
	require.Equal(t, "aac", i.Codec)
	require.Equal(t, 64.0, i.Bitrate)
	require.Equal(t, 0.0, i.Duration)
	require.Equal(t, uint64(44100), i.Sampling)
	require.Equal(t, "stereo", i.Layout)
}
