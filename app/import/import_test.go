package main

import (
	gojson "encoding/json"
	"os"
	"testing"

	"github.com/datarhei/core/v16/encoding/json"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/restream/store"

	"github.com/stretchr/testify/require"
)

// Scenarios:
// - Empty v1 file - v1/v4 empty
// - Input: RTSP (h264/), copy video, no audio, Output: none - v1/v4 rtsp_h264,-_copyvideo_noaudio_output_none
// - Input: HLS (h264/), copy video, no audio, Output: none - v1/v4 hls_h264,-_copyvideo_noaudio_output_none
// - Input: RTMP (h264/), copy video, no audio, Output: none - v1/v4 rtmp_h264,-_copyvideo_noaudio_output_none
// - Input: RTSP (h264/), copy video, auto audio, Output: none - v1/v4 rtsp_h264,-_copyvideo_autoaudio_output_none
// - Input: RTSP (h264/aac), copy video, copy audio, Output: none - v1/v4 rtsp_h264,aac_copyvideo_copyaudio_output_none
// - Input: RTSP (h264/mp3), copy video, copy audio, Output: none - v1/v4 rtsp_h264,mp3_copyvideo_copyaudio_output_none
// - Input: RTSP (h264/aac), copy video, encode audio aac, Output: none - v1/v4 rtsp_h264,aac_copyvideo_encodeaudioaac_output_none
// - Input: RTSP (h264/mp3), copy video, encode audio mp3, Output: none - v1/v4 rtsp_h264,mp3_copyvideo_encodeaudiomp3_output_none
// - Input: RTSP (h264/aac), copy video, no audio, Output: none - v1/v4 rtsp_h264,aac_copyvideo_noaudio_output_none
// - Input: RTSP (h264/aac), copy video, encode silence aac, Output: none - v1/v4 rtsp_h264,aac_copyvideo_encodesilenceaac_output_none
// - Input: RTSP (h264/), copy video, encode silence aac, Output: none - v1/v4 rtsp_h264,-_copyvideo_encodesilenceaac_output_none
// - Input: RTSP (h264/), copy video, encode silence mp3, Output: none - v1/v4 rtsp_h264_copyvideo_encodesilencemp3_output_none
// - Input: RTSP (h264/), encode video, no audio, Output: none - v1v4/ rtsp_h264,-_encodevideo_noaudio_output_none
// - Input: RTSP (h264/aac), encode video, copy audio, Output: none - v1/v4 rtsp_h264,aac_encodevideo_copyaudio_output_none
// - Input: RTSP (h264/aac), encode video, encode audio aac, Output: none - v1/v4 rtsp_h264,aac_encodevideo_encodeaudioaac_output_none
// - Input: RTSP (h264/aac), encode video, encode silence aac, Output: none - v1/v4 rtsp_h264,aac_encodevideo_encodesilenceaac_output_none
// - Input: RTSP (h264/aac), encode video, no audio, Output: none - v1/v4 rtsp_h264,aac_encodevideo_noaudio_output_none
// - Input: RTSP (h264/aac), copy video, copy audio, Output: RTMP - v1/v4 rtsp_h264,aac_copyvideo_copyaudio_output_rtmp
// - Input: RTSP (h264/aac), copy video, copy audio, Output: HLS - v1/v4 rtsp_h264,aac_copyvideo_copyaudio_output_hls

var id string = "4186b095-7f0a-4e94-8c3d-f17459ab252f"

func testV1Import(t *testing.T, v1Fixture, v4Fixture string, config importConfig) {
	diskfs, err := fs.NewRootedDiskFilesystem(fs.RootedDiskConfig{
		Root: ".",
	})
	require.NoError(t, err)

	// Import v1 database
	v4, err := importV1(diskfs, v1Fixture, config)
	require.Equal(t, nil, err)

	// Reset variants
	for n := range v4.Process {
		v4.Process[n].CreatedAt = 0
	}

	// Convert to JSON
	datav4, err := gojson.MarshalIndent(&v4, "", "    ")
	require.Equal(t, nil, err)

	// Read the wanted result
	wantdatav4, err := diskfs.ReadFile(v4Fixture)
	require.Equal(t, nil, err)

	var wantv4 store.StoreData

	err = gojson.Unmarshal(wantdatav4, &wantv4)
	require.Equal(t, nil, err, json.FormatError(wantdatav4, err))

	// Convert to JSON
	wantdatav4, err = gojson.MarshalIndent(&wantv4, "", "    ")
	require.Equal(t, nil, err)

	// Re-convert both to golang type
	gojson.Unmarshal(wantdatav4, &wantv4)
	gojson.Unmarshal(datav4, &v4)

	require.Equal(t, wantv4, v4)
}

func TestV1Import(t *testing.T) {
	tests := []string{
		"empty",
		"hls_h264,-_copyvideo_noaudio_output_none",
		"rtmp_h264,-_copyvideo_noaudio_output_none",
		"rtsp_h264,-_copyvideo_autoaudio_output_none",
		"rtsp_h264,-_copyvideo_encodesilenceaac_output_none",
		"rtsp_h264,-_copyvideo_encodesilencemp3_output_none",
		"rtsp_h264,-_copyvideo_noaudio_output_none",
		"rtsp_h264,-_encodevideo_noaudio_output_none",
		"rtsp_h264,aac_copyvideo_copyaudio_output_hls",
		"rtsp_h264,aac_copyvideo_copyaudio_output_none",
		"rtsp_h264,aac_copyvideo_copyaudio_output_rtmp",
		"rtsp_h264,aac_copyvideo_encodeaudioaac_output_none",
		"rtsp_h264,aac_copyvideo_encodesilenceaac_output_none",
		"rtsp_h264,aac_copyvideo_noaudio_output_none",
		"rtsp_h264,aac_encodevideo_copyaudio_output_none",
		"rtsp_h264,aac_encodevideo_encodeaudioaac_output_none",
		"rtsp_h264,aac_encodevideo_encodesilenceaac_output_none",
		"rtsp_h264,aac_encodevideo_noaudio_output_none",
		"rtsp_h264,mp3_copyvideo_copyaudio_output_none",
		"rtsp_h264,mp3_copyvideo_encodeaudiomp3_output_none",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			testV1Import(t, "./fixtures/v1_"+test+".json", "./fixtures/v4_"+test+".json", importConfig{
				id:               id,
				snapshotInterval: 60,
			})
		})
	}
}

func TestImportSnapshotInterval(t *testing.T) {
	type testdata struct {
		value    string
		expected int
	}

	tests := []testdata{
		{"0", 0},
		{"0ms", 0},
		{"0s", 0},
		{"0m", 0},
		{"60000", 60},
		{"60000ms", 60},
		{"60s", 60},
		{"1m", 60},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			actual := importSnapshotInterval(test.value, -1)
			require.Equal(t, test.expected, actual)
		})
	}
}

// Scenarios:
// Input: V4L (h264), copy video, auto audio, Output: none - v1/v4 v4l_h264,-_copyvideo_autoaudio_output_none
// Input: V4L (h264), copy video, no audio, Output: none - v1/v4 v4l_h264,-_copyvideo_noaudio_output_none
// Input: V4L (h264), copy video, encode silence aac, Output: none - v1/v4 v4l_h264,-_copyvideo_encodesilenceaac_output_none
// Input: V4L (h264), encode video, auto audio, Output: none - v1/v4 v4l_h264,-_encodevideo_autoaudio_output_none

func TestImportUSBCamWithoutAudio(t *testing.T) {
	tests := []string{
		"v4l_h264,-_copyvideo_autoaudio_output_none",
		"v4l_h264,-_copyvideo_noaudio_output_none",
		"v4l_h264,-_copyvideo_encodesilenceaac_output_none",
		"v4l_h264,-_encodevideo_autoaudio_output_none",
	}

	config := importConfig{
		id:               id,
		snapshotInterval: 60,
		usbcam: importConfigUSBCam{
			enable:  true,
			device:  "/dev/video",
			fps:     "25",
			gop:     "50",
			bitrate: "5000000",
			preset:  "ultrafast",
			profile: "auto",
			width:   "1280",
			height:  "720",
		},
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			testV1Import(t, "./fixtures/v1_"+test+".json", "./fixtures/v4_"+test+".json", config)
		})
	}
}

// Scenarios:
// Input: V4L/ALSA (h264/aac), copy video, auto audio, Output: none - v1/v4 v4lalsa_h264,aac_copyvideo_autoaudio_output_none
// Input: V4L/ALSA (h264/aac), copy video, copy audio, Output: none - v1/v4 v4lalsa_h264,aac_copyvideo_copyaudio_output_none
// Input: V4L/ALSA (h264/aac), copy video, no audio, Output: none - v1/v4 v4lalsa_h264,aac_copyvideo_noaudio_output_none
// Input: V4L/ALSA (h264/aac), copy video, encode silence aac, Output: none - v1/v4 v4lalsa_h264,aac_copyvideo_encodesilenceaac_output_none
// Input: V4L/ALSA (h264/aac), copy video, encode audio aac, Output: none - v1/v4 v4lalsa_h264,aac_copyvideo_encodeaudioaac_output_none
// Input: V4L/ALSA (h264/aac), encode video, auto audio, Output: none - v1/v4 v4lalsa_h264,aac_encodevideo_autoaudio_output_none

func TestImportUSBCamWithAudio(t *testing.T) {
	tests := []string{
		"v4lalsa_h264,aac_copyvideo_autoaudio_output_none",
		"v4lalsa_h264,aac_copyvideo_copyaudio_output_none",
		"v4lalsa_h264,aac_copyvideo_noaudio_output_none",
		"v4lalsa_h264,aac_copyvideo_encodesilenceaac_output_none",
		"v4lalsa_h264,aac_copyvideo_encodeaudioaac_output_none",
		"v4lalsa_h264,aac_encodevideo_autoaudio_output_none",
	}

	config := importConfig{
		id:               id,
		snapshotInterval: 60,
		usbcam: importConfigUSBCam{
			enable:  true,
			device:  "/dev/video",
			fps:     "25",
			gop:     "50",
			bitrate: "5000000",
			preset:  "ultrafast",
			profile: "auto",
			width:   "1280",
			height:  "720",
		},
		audio: importConfigAudio{
			enable:   true,
			device:   "1,0",
			bitrate:  "28000",
			channels: "2",
			layout:   "stereo",
			sampling: "11000",
		},
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			testV1Import(t, "./fixtures/v1_"+test+".json", "./fixtures/v4_"+test+".json", config)
		})
	}
}

func TestV1EnvironmentDefaults(t *testing.T) {
	for key := range v1Environment {
		os.Unsetenv(key)
	}

	initV1Environment()

	for _, val := range v1Environment {
		require.Equal(t, val.defval, val.value)
	}

	config := importConfigFromEnvironment()

	require.Equal(t, false, config.usbcam.enable)
	require.Equal(t, false, config.raspicam.enable)
	require.Equal(t, false, config.audio.enable)
}

func TestV1EnvironmentUSBCam(t *testing.T) {
	for key := range v1Environment {
		os.Unsetenv(key)
	}

	os.Setenv("RS_MODE", "USBCAM")

	config := importConfigFromEnvironment()

	require.Equal(t, true, config.usbcam.enable)
	require.Equal(t, "/dev/video", config.usbcam.device)
	require.Equal(t, "25", config.usbcam.fps)
	require.Equal(t, "50", config.usbcam.gop)
	require.Equal(t, "5000000", config.usbcam.bitrate)
	require.Equal(t, "ultrafast", config.usbcam.preset)
	require.Equal(t, "baseline", config.usbcam.profile)
	require.Equal(t, "1280", config.usbcam.width)
	require.Equal(t, "720", config.usbcam.height)

	os.Setenv("RS_USBCAM_AUDIO", "true")

	config = importConfigFromEnvironment()

	require.Equal(t, true, config.audio.enable)
	require.Equal(t, "0", config.audio.device)
	require.Equal(t, "64000", config.audio.bitrate)
	require.Equal(t, "1", config.audio.channels)
	require.Equal(t, "mono", config.audio.layout)
	require.Equal(t, "44100", config.audio.sampling)
}

func TestV1EnvironmentRASPICam(t *testing.T) {
	for key := range v1Environment {
		os.Unsetenv(key)
	}

	os.Setenv("RS_MODE", "RASPICAM")

	config := importConfigFromEnvironment()

	require.Equal(t, true, config.raspicam.enable)
	require.Equal(t, "25", config.raspicam.fps)
	require.Equal(t, "1920", config.raspicam.width)
	require.Equal(t, "1080", config.raspicam.height)

	os.Setenv("RS_RASPICAM_AUDIO", "true")

	config = importConfigFromEnvironment()

	require.Equal(t, true, config.audio.enable)
	require.Equal(t, "0", config.audio.device)
	require.Equal(t, "64000", config.audio.bitrate)
	require.Equal(t, "1", config.audio.channels)
	require.Equal(t, "mono", config.audio.layout)
	require.Equal(t, "44100", config.audio.sampling)
}

func TestV1EnvironmentInputstream(t *testing.T) {
	for key := range v1Environment {
		os.Unsetenv(key)
	}

	os.Setenv("RS_INPUTSTREAM", "")
}

func TestV1Pre067(t *testing.T) {
	tests := []string{
		"pre-0.6.7",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			testV1Import(t, "./fixtures/v1_"+test+".json", "./fixtures/v4_"+test+".json", importConfig{
				id:               id,
				snapshotInterval: 60,
				binary:           "ffmpeg",
			})
		})
	}
}
