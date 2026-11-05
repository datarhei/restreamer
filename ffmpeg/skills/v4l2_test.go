package skills

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoV4LDevices(t *testing.T) {
	data := bytes.NewBufferString(``)

	devices := parseV4LDevices(data)

	require.Equal(t, []HWDevice{}, devices)
}

func TestV4LDevices(t *testing.T) {
	data := bytes.NewBufferString(`mmal service 16.1 (platform:bcm2835-v4l2):
	/dev/video0

Webcam C170: Webcam C170 (usb-3f980000.usb-1.3):
	/dev/video1

`)

	devices := parseV4LDevices(data)

	require.Equal(t, []HWDevice{
		{
			Id:    "/dev/video0",
			Name:  "mmal service 16.1",
			Extra: "platform:bcm2835-v4l2",
			Media: "video",
		},
		{
			Id:    "/dev/video1",
			Name:  "Webcam C170: Webcam C170",
			Extra: "usb-3f980000.usb-1.3",
			Media: "video",
		},
	}, devices)
}
