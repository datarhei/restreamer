package skills

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoALSACards(t *testing.T) {
	data := ``

	cards := parseALSACards([]byte(data))

	require.Equal(t, []alsaCard{}, cards)
}

func TestALSACards(t *testing.T) {
	data := ` 0 [ALSA           ]: bcm2835_alsa - bcm2835 ALSA
	bcm2835 ALSA
1 [C170           ]: USB-Audio - Webcam C170
	Webcam C170 at usb-3f980000.usb-1.3, high speed`

	cards := parseALSACards([]byte(data))

	require.Equal(t, []alsaCard{
		{
			id:   "0",
			name: "bcm2835_alsa - bcm2835 ALSA",
		},
		{
			id:   "1",
			name: "USB-Audio - Webcam C170",
		},
	}, cards)
}

func TestNoALSADevices(t *testing.T) {
	data := ``

	cards := []alsaCard{}

	devices := parseALSADevices([]byte(data), cards)

	require.Equal(t, []HWDevice{}, devices)
}

func TestALSADevices(t *testing.T) {
	data := `  0: [ 0]   : control
	16: [ 0- 0]: digital audio playback
	17: [ 0- 1]: digital audio playback
	32: [ 1]   : control
	33:        : timer
	56: [ 1- 0]: digital audio capture`

	cards := []alsaCard{
		{
			id:   "0",
			name: "bcm2835_alsa - bcm2835 ALSA",
		},
		{
			id:   "1",
			name: "USB-Audio - Webcam C170",
		},
	}

	devices := parseALSADevices([]byte(data), cards)

	require.Equal(t, []HWDevice{
		{
			Id:    "hw:1,0",
			Name:  "USB-Audio - Webcam C170",
			Extra: "",
			Media: "audio",
		},
	}, devices)
}
