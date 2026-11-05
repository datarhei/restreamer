package skills

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"strings"
)

// DevicesV4L returns a list of available V4L devices
func DevicesV4L() ([]HWDevice, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command("v4l2-ctl", "--list-devices")
	cmd.Env = []string{}
	cmd.Stdout = buf
	cmd.Run()

	devices := parseV4LDevices(buf)

	return devices, nil
}

func parseV4LDevices(data *bytes.Buffer) []HWDevice {
	devices := []HWDevice{}

	reName := regexp.MustCompile(`^(.*)\((.*?)\):$`)
	reDevice := regexp.MustCompile(`/dev/video[0-9]+`)

	name := ""
	extra := ""

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := reDevice.FindStringSubmatch(line)
		if matches == nil {
			matches = reName.FindStringSubmatch(line)
			if matches == nil {
				continue
			}

			name = strings.TrimSpace(matches[1])
			extra = matches[2]

			continue
		}

		if name == "" {
			continue
		}

		device := HWDevice{
			Id:    matches[0],
			Name:  name,
			Extra: extra,
			Media: "video",
		}

		devices = append(devices, device)
	}

	return devices
}
