package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	header := `ffmpeg version 4.0.2 Copyright (c) 2000-2018 the FFmpeg developers
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
  libpostproc    55.  1.100 / 55.  1.100`

	prelude := `Input #0, lavfi, from 'testsrc=size=1280x720:rate=25':
  Duration: N/A, start: 0.000000, bitrate: N/A
    Stream #0:0: Video: rawvideo (RGB[24] / 0x18424752), rgb24, 1280x720 [SAR 1:1 DAR 16:9], 25 tbr, 25 tbn, 25 tbc
Input #1, lavfi, from 'anullsrc=r=44100:cl=stereo':
  Duration: N/A, start: 0.000000, bitrate: 705 kb/s
    Stream #1:0: Audio: pcm_u8, 44100 Hz, stereo, u8, 705 kb/s
	Stream #1:1(eng): Audio: aac (LC), 48000 Hz, stereo, fltp (default)
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
[hls @ 0x7fa969803a00] Opening './data/testsrc.m3u8.tmp' for writing`

	fmt.Fprintf(os.Stderr, "%s\n", header)

	if len(os.Args) <= 1 {
		os.Exit(2)
	}

	lastArg := os.Args[len(os.Args)-1]

	if lastArg == "-version" {
		os.Exit(0)
	}

	if len(lastArg) > 1 && lastArg[0] == '-' {
		os.Exit(2)
	}

	fmt.Fprintf(os.Stderr, "%s\n", prelude)

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		frame := uint64(0)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				frame += 25
				fmt.Fprintf(os.Stderr, "frame=%5d fps= 25 q=-1.0 Lsize=N/A time=00:00:02.32 bitrate=N/A speed=1.0x    \r", frame)
			}
		}
	}(ctx)

	// Wait for interrupt signal to gracefully shutdown the app
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cancel()

	fmt.Fprintf(os.Stderr, "\nExiting normally, received signal 2.\n")

	os.Exit(255)
}
