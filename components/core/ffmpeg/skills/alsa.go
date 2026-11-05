package skills

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
)

type alsaCard struct {
	id   string
	name string
}

// DevicesALSA returns a list of available ALSA devices
func DevicesALSA() ([]HWDevice, error) {
	devices := []HWDevice{}

	content, err := os.ReadFile("/proc/asound/cards")
	if err != nil {
		return devices, err
	}

	cards := parseALSACards(content)

	content, err = os.ReadFile("/proc/asound/devices")
	if err != nil {
		return devices, err
	}

	devices = parseALSADevices(content, cards)

	return devices, nil
}

func parseALSACards(data []byte) []alsaCard {
	cards := []alsaCard{}

	reCards := regexp.MustCompile(`^\s*([0-9]+) \[([^\]]+)\]: (.*)$`)

	buf := bytes.NewBuffer(data)

	scanner := bufio.NewScanner(buf)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := reCards.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		cards = append(cards, alsaCard{
			id:   matches[1],
			name: matches[3],
		})
	}

	return cards
}

func parseALSADevices(data []byte, cards []alsaCard) []HWDevice {
	devices := []HWDevice{}

	reDevices := regexp.MustCompile(`^\s*[0-9]+: \[\s?([0-9]+)-\s?([0-9]+)\]: (.*)$`)

	buf := bytes.NewBuffer(data)

	scanner := bufio.NewScanner(buf)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()

		matches := reDevices.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		if matches[3] != "digital audio capture" {
			continue
		}

		// find the card
		for _, card := range cards {
			if card.id != matches[1] {
				continue
			}

			device := HWDevice{
				Id:    "hw:" + matches[1] + "," + matches[2],
				Name:  card.name,
				Media: "audio",
			}

			devices = append(devices, device)

			break
		}
	}

	return devices
}
