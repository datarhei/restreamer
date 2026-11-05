package store

import (
	"os"
	"path"
)

// Location returns the path to the config file. If no path is provided,
// different standard location will be probed:
// - os.UserConfigDir() + /datarhei-core/config.json
// - os.UserHomeDir() + /.config/datarhei-core/config.json
// - ./config/config.json
// If the config doesn't exist in none of these locations, it will be assumed
// at ./config/config.json
func Location(filepath string) string {
	configfile := filepath
	if len(configfile) != 0 {
		return configfile
	}

	locations := []string{}

	if dir, err := os.UserConfigDir(); err == nil {
		locations = append(locations, dir+"/datarhei-core/config.js")
		locations = append(locations, dir+"/datarhei-core/config.json")
	}

	if dir, err := os.UserHomeDir(); err == nil {
		locations = append(locations, dir+"/.config/datarhei-core/config.js")
		locations = append(locations, dir+"/.config/datarhei-core/config.json")
	}

	locations = append(locations, "./config/config.js")
	locations = append(locations, "./config/config.json")

	for _, path := range locations {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}

		if info.IsDir() {
			continue
		}

		configfile = path
	}

	if len(configfile) == 0 {
		configfile = "./config/config.json"
	}

	os.MkdirAll(path.Dir(configfile), 0740)

	return configfile
}
