package main

import (
	"fmt"
	"os"

	cfgstore "github.com/datarhei/core/v16/config/store"
	cfgvars "github.com/datarhei/core/v16/config/vars"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/restream/store"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := log.New("Import").WithOutput(log.NewConsoleWriter(os.Stderr, log.Linfo, true)).WithField("version", "v1")

	configfile := cfgstore.Location(os.Getenv("CORE_CONFIGFILE"))

	diskfs, err := fs.NewDiskFilesystem(fs.DiskConfig{})
	if err != nil {
		logger.Error().WithError(err).Log("Access disk filesystem failed")
		os.Exit(1)
	}

	configstore, err := cfgstore.NewJSON(diskfs, configfile, nil)
	if err != nil {
		logger.Error().WithError(err).Log("Loading configuration failed")
		os.Exit(1)
	}

	if err := doImport(logger, diskfs, configstore); err != nil {
		os.Exit(1)
	}
}

func doImport(logger log.Logger, fs fs.Filesystem, configstore cfgstore.Store) error {
	if logger == nil {
		logger = log.New("")
	}

	logger.Info().Log("Database import")

	cfg := configstore.Get()

	// Merging the persisted config with the environment variables
	cfg.Merge()

	cfg.Validate(false)
	if cfg.HasErrors() {
		logger.Error().Log("The configuration contains errors")
		messages := []string{}
		cfg.Messages(func(level string, v cfgvars.Variable, message string) {
			if level == "error" {
				logger.Error().WithFields(log.Fields{
					"variable":    v.Name,
					"value":       v.Value,
					"env":         v.EnvName,
					"description": v.Description,
				}).Log(message)

				messages = append(messages, v.Name+": "+message)
			}
		})

		return fmt.Errorf("the configuration contains errors: %v", messages)
	}

	logger.Info().Log("Checking for database ...")

	// Check if there's a v1.json from the old Restreamer
	v1filename := cfg.DB.Dir + "/v1.json"

	logger = logger.WithField("database", v1filename)

	if _, err := fs.Stat(v1filename); err != nil {
		if os.IsNotExist(err) {
			logger.Info().Log("Database doesn't exist and nothing will be imported")
			return nil
		}

		logger.Error().WithError(err).Log("Checking for v1 database")
		return fmt.Errorf("checking for v1 database: %w", err)
	}

	logger.Info().Log("Found database")

	// Load an existing DB
	datastore, err := store.NewJSON(store.JSONConfig{
		Filesystem: fs,
		Filepath:   cfg.DB.Dir + "/db.json",
	})
	if err != nil {
		logger.Error().WithError(err).Log("Creating datastore for new database failed")
		return fmt.Errorf("creating datastore for new database failed: %w", err)
	}

	data, err := datastore.Load()
	if err != nil {
		logger.Error().WithError(err).Log("Loading new database failed")
		return fmt.Errorf("loading new database failed: %w", err)
	}

	// Check if the existing DB has already some data in it.
	// If it's not empty, we will not import any v1 DB.
	if !data.IsEmpty() {
		logger.Info().Log("There's already information stored in the new database and the v1 database will not be imported")
		return nil
	}

	logger.Info().Log("Importing database ...")

	// Read the Restreamer config from the environment variables
	importConfig := importConfigFromEnvironment()
	importConfig.binary = cfg.FFmpeg.Binary

	// Rewrite the old database to the new database
	r, err := importV1(fs, v1filename, importConfig)
	if err != nil {
		logger.Error().WithError(err).Log("Importing database failed")
		return fmt.Errorf("importing database failed: %w", err)
	}

	// Persist the imported DB
	if err := datastore.Store(r); err != nil {
		logger.Error().WithError(err).Log("Storing imported data to new database failed")
		return fmt.Errorf("storing imported data to new database failed: %w", err)
	}

	// Get the unmerged config for persisting
	cfg = configstore.Get()

	// Add static routes to mimic the old URLs
	cfg.Router.Routes["/hls/live.stream.m3u8"] = "/memfs/" + importConfig.id + ".m3u8"
	cfg.Router.Routes["/images/live.jpg"] = "/memfs/" + importConfig.id + ".jpg"
	cfg.Router.Routes["/player.html"] = "/" + importConfig.id + ".html"

	// Persist the modified config
	if err := configstore.Set(cfg); err != nil {
		logger.Error().WithError(err).Log("Storing adjusted config failed")
		return fmt.Errorf("storing adjusted config failed: %w", err)
	}

	logger.Info().Log("Successfully imported data")

	return nil
}
