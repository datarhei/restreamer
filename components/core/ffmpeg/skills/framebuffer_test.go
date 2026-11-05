package skills

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoFramebufferDevices(t *testing.T) {
	data := bytes.NewBufferString(``)

	extra := parseFramebufferDevice(data)

	require.Equal(t, "", extra)
}

func TestFramebufferDevices(t *testing.T) {
	data := []*bytes.Buffer{}

	data = append(data, bytes.NewBufferString(`mode "1280x720"
    geometry 1280 720 1280 720 32
    timings 0 0 0 0 0 0 0
    rgba 8/16,8/8,8/0,8/24
endmode`))

	data = append(data, bytes.NewBufferString(`mode "1280x720"
    geometry 1280 720 1280 720 16
    timings 0 0 0 0 0 0 0
    rgba 5/11,6/5,5/0,0/16
endmode`))

	data = append(data, bytes.NewBufferString(`mode "1280x720"
    geometry 1280 720 1280 720 8
    timings 0 0 0 0 0 0 0
    rgba 8/0,8/0,8/0,0/0
endmode`))

	data = append(data, bytes.NewBufferString(`mode "1280x720"
    geometry 1280 720 1280 720 24
    timings 0 0 0 0 0 0 0
    rgba 8/16,8/8,8/0,0/24
endmode`))

	extras := []string{}

	for _, d := range data {
		extras = append(extras, parseFramebufferDevice(d))
	}

	require.Equal(t, []string{
		"1280x720 bgra",
		"1280x720 rgb565le",
		"",
		"1280x720 bgr24",
	}, extras)
}
