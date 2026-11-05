package skills

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"strings"
)

// DevicesAvfoundation returns a list of AVFoundation devices. You have to provide
// the path to the ffmpeg binary.
func DevicesAvfoundation(binary string) ([]HWDevice, error) {
	buf := bytes.NewBuffer(nil)

	cmd := exec.Command(binary, "-f", "avfoundation", "-list_devices", "true", "-i", "")
	cmd.Env = []string{}
	cmd.Stderr = buf
	cmd.Run()

	devices := parseAvfoundationDevices(buf)

	return devices, nil
}

func parseAvfoundationDevices(data *bytes.Buffer) []HWDevice {
	devices := []HWDevice{}

	re := regexp.MustCompile(`^\s*\[.*?\] \[([0-9]+)\] (.*?)$`)

	media := ""

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "AVFoundation video devices") {
			media = "video"
			continue
		} else if strings.Contains(line, "AVFoundation audio devices") {
			media = "audio"
			continue
		}

		if media == "" {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		device := HWDevice{
			Id:    matches[1],
			Name:  matches[2],
			Media: media,
		}

		devices = append(devices, device)
	}

	return devices
}
