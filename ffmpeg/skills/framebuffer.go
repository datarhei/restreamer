package skills

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type fbformat struct {
	bpp     int
	roffset int
	goffset int
	boffset int
	aoffset int
	format  string
}

var fbformats []fbformat = []fbformat{
	{32, 0, 8, 16, 24, "rgba"},
	{32, 16, 8, 0, 24, "bgra"},
	{32, 8, 16, 24, 0, "argb"},
	{32, 3, 2, 8, 0, "abgr"},
	{24, 0, 8, 16, 0, "rgb24"},
	{24, 16, 8, 0, 0, "bgr24"},
	{16, 11, 5, 0, 0, "rgb565le"},
	{16, 0, 5, 11, 0, "rgb565be"},
}

// DevicesFramebuffer returns a list of framebuffer devices found on the system
func DevicesFramebuffer() ([]HWDevice, error) {
	devices := []HWDevice{}

	for i := 0; i < 32; i++ {
		path := fmt.Sprintf("/dev/fb%d", i)
		_, err := os.Stat(path)
		if err != nil {
			break
		}

		buf := bytes.NewBuffer(nil)

		cmd := exec.Command("fbset", "-s", "-fb", path)
		cmd.Env = []string{}
		cmd.Stdout = buf
		cmd.Run()

		device := HWDevice{
			Id:    path,
			Name:  path,
			Extra: parseFramebufferDevice(buf),
			Media: "video",
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func parseFramebufferDevice(data *bytes.Buffer) string {
	reGeometry := regexp.MustCompile(`geometry ([0-9]+) ([0-9]+) ([0-9]+) ([0-9]+) ([0-9]+)`)
	reRGBA := regexp.MustCompile(`rgba ([0-9]+/[0-9]+),([0-9]+/[0-9]+),([0-9]+/[0-9]+),([0-9]+/[0-9]+)`)

	geometry := ""
	format := fbformat{}

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := reGeometry.FindStringSubmatch(line)
		if matches != nil {
			geometry = fmt.Sprintf("%sx%s", matches[1], matches[2])
			if matches[1] != matches[3] || matches[2] != matches[4] {
				geometry += fmt.Sprintf(" (%sx%s)", matches[3], matches[4])
			}

			if x, err := strconv.Atoi(matches[5]); err == nil {
				format.bpp = x
			}
		}

		matches = reRGBA.FindStringSubmatch(line)
		if matches != nil {
			format.roffset = parseFramebufferDeviceOffset(matches[1])
			format.goffset = parseFramebufferDeviceOffset(matches[2])
			format.boffset = parseFramebufferDeviceOffset(matches[3])
			format.aoffset = parseFramebufferDeviceOffset(matches[4])
		}
	}

	if len(geometry) == 0 || format.bpp == 0 {
		return ""
	}

	for _, f := range fbformats {
		if f.bpp == format.bpp && f.roffset == format.roffset && f.goffset == format.goffset && f.boffset == format.boffset {
			return geometry + " " + f.format
		}
	}

	return ""
}

func parseFramebufferDeviceOffset(format string) int {
	c := strings.Split(format, "/")

	bits := 0
	offset := 0

	if x, err := strconv.Atoi(c[0]); err == nil {
		bits = x
	} else {
		return 0
	}

	if x, err := strconv.Atoi(c[1]); err == nil {
		offset = x
	} else {
		return 0
	}

	if bits == 0 {
		return 0
	}

	return offset
}
