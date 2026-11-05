package skills

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoAvfoundationDevices(t *testing.T) {
	data := bytes.NewBufferString(``)

	devices := parseAvfoundationDevices(data)

	require.Equal(t, []HWDevice{}, devices)
}

func TestAvfoundationDevices(t *testing.T) {
	data := bytes.NewBufferString(`[AVFoundation input device @ 0x7fc2db40f240] AVFoundation video devices:
	[AVFoundation input device @ 0x7fc2db40f240] [0] FaceTime HD Camera (Built-in)
	[AVFoundation input device @ 0x7fc2db40f240] [1] Capture screen 0
	[AVFoundation input device @ 0x7fc2db40f240] [2] Capture screen 1
	[AVFoundation input device @ 0x7fc2db40f240] AVFoundation audio devices:
	[AVFoundation input device @ 0x7fc2db40f240] [0] Built-in Microphone
	: Input/output error`)

	devices := parseAvfoundationDevices(data)

	require.Equal(t, []HWDevice{
		{
			Id:    "0",
			Name:  "FaceTime HD Camera (Built-in)",
			Extra: "",
			Media: "video",
		},
		{
			Id:    "1",
			Name:  "Capture screen 0",
			Extra: "",
			Media: "video",
		},
		{
			Id:    "2",
			Name:  "Capture screen 1",
			Extra: "",
			Media: "video",
		},
		{
			Id:    "0",
			Name:  "Built-in Microphone",
			Extra: "",
			Media: "audio",
		},
	}, devices)
}
