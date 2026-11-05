package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	cfgstore "github.com/datarhei/core/v16/config/store"
	cfgvars "github.com/datarhei/core/v16/config/vars"
	"github.com/datarhei/core/v16/ffmpeg"
	"github.com/datarhei/core/v16/io/file"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/log"
	"github.com/datarhei/core/v16/restream/app"
	"github.com/datarhei/core/v16/restream/store"

	"github.com/Masterminds/semver/v3"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := log.New("Migration").WithOutput(log.NewConsoleWriter(os.Stderr, log.Linfo, true))

	configfile := cfgstore.Location(os.Getenv("CORE_CONFIGFILE"))

	diskfs, _ := fs.NewDiskFilesystem(fs.DiskConfig{})

	configstore, err := cfgstore.NewJSON(diskfs, configfile, nil)
	if err != nil {
		logger.Error().WithError(err).Log("Loading configuration failed")
		os.Exit(1)
	}

	if err := doMigration(logger, diskfs, configstore); err != nil {
		logger.Error().WithError(err).Log("Migration failed")
		os.Exit(1)
	}
}

func doMigration(logger log.Logger, fs fs.Filesystem, configstore cfgstore.Store) error {
	if logger == nil {
		logger = log.New("")
	}

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

	ff, err := ffmpeg.New(ffmpeg.Config{
		Binary: cfg.FFmpeg.Binary,
	})
	if err != nil {
		logger.Error().WithError(err).Log("Loading FFmpeg binary failed")
		return fmt.Errorf("loading FFmpeg binary failed: %w", err)
	}

	version, err := semver.NewVersion(ff.Skills().FFmpeg.Version)
	if err != nil {
		logger.Error().WithError(err).Log("Parsing FFmpeg version failed")
		return fmt.Errorf("parsing FFmpeg version failed: %w", err)
	}

	// Check if there's a DB file
	dbFilepath := cfg.DB.Dir + "/db.json"

	if _, err = os.Stat(dbFilepath); err != nil {
		// There's no DB to backup
		logger.Info().WithField("db", dbFilepath).Log("Database not found. Migration not required")
		return nil
	}

	// Load the existing DB
	datastore, err := store.NewJSON(store.JSONConfig{
		Filesystem: fs,
		Filepath:   cfg.DB.Dir + "/db.json",
	})
	if err != nil {
		logger.Error().WithField("db", dbFilepath).WithError(err).Log("Creating JSON store failed")
		return fmt.Errorf("creating JSON store failed: %w", err)
	}

	data, err := datastore.Load()
	if err != nil {
		logger.Error().WithError(err).Log("Loading database failed")
		return fmt.Errorf("loading database failed: %w", err)
	}

	// Migrate processes
	logger.Info().Log("Migrating processes ...")

	migrated := false

	for id, p := range data.Process {
		ok, err := migrateProcessConfig(logger.WithField("processid", p.ID), p.Config, version.String())
		if err != nil {
			logger.Info().WithField("processid", p.ID).WithError(err).Log("Migrating process failed")
			return fmt.Errorf("migrating process failed: %w", err)
		}

		data.Process[id] = p

		if ok {
			migrated = true
		}
	}

	logger.Info().Log("Migrating processes done")

	if migrated {
		// Create backup if something has been changed.
		backupFilepath := cfg.DB.Dir + "/db." + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".json"

		if _, err = os.Stat(backupFilepath); err == nil {
			// Yes, we have a backup. The migration already happened
			logger.Info().WithField("backup", backupFilepath).Log("Migration already done")
			return nil
		}

		// Create a backup
		if err := file.Copy(dbFilepath, backupFilepath); err != nil {
			logger.Error().WithError(err).Log("Creating backup file failed")
			return fmt.Errorf("creating backup file failed: %w", err)
		}

		logger.Info().WithField("backup", backupFilepath).Log("Backup created")

		// Store the modified DB
		if err := datastore.Store(data); err != nil {
			logger.Error().WithError(err).Log("Storing database failed")
			return fmt.Errorf("storing database failed: %w", err)
		}
	}

	logger.Info().Log("Completed")

	return nil
}

func migrateProcessConfig(logger log.Logger, config *app.Config, version string) (bool, error) {
	migrated := false

	vtarget, err := semver.NewVersion(version)
	if err != nil {
		logger.Error().WithError(err).Log("Parsing target FFmpeg version failed")
		return false, fmt.Errorf("parsing target FFmpeg version failed: %w", err)
	}

	targetmajor := vtarget.Major()
	currentmajor := uint64(4)

	if len(config.FFVersion) != 0 {
		vcurrent, err := semver.NewVersion(strings.TrimPrefix(config.FFVersion, "^"))
		if err != nil {
			logger.Error().WithError(err).Log("Parsing current FFmpeg version failed")
			return false, fmt.Errorf("parsing current FFmpeg version failed: %w", err)
		}

		currentmajor = vcurrent.Major()
	}

	if currentmajor < 4 {
		err := fmt.Errorf("unknown FFmpeg version found: %d", currentmajor)
		logger.Error().WithError(err).Log("Unsupported FFmpeg version")
		return false, fmt.Errorf("unsupported FFmpeg version: %w", err)
	}

	if targetmajor > 6 {
		err := fmt.Errorf("unknown FFmpeg version found: %d", targetmajor)
		logger.Error().WithError(err).Log("Unsupported FFmpeg version")
		return false, fmt.Errorf("unsupported FFmpeg version: %w", err)
	}

	if currentmajor != targetmajor {
		migrated = true
	}

	if currentmajor == 4 && targetmajor > 4 {
		// Migration from version 4 to version 5
		// Only this happens:
		// - for RTSP inputs, replace -stimeout with -timeout
		reRTSP := regexp.MustCompile(`^rtsps?://`)

		for index, input := range config.Input {
			if !reRTSP.MatchString(input.Address) {
				continue
			}

			for i, o := range input.Options {
				if o != "-stimeout" {
					continue
				}

				input.Options[i] = "-timeout"
			}

			config.Input[index] = input
		}

		currentmajor = 5
	}

	if currentmajor == 5 && targetmajor > 5 {
		// Migration from version 5 to version 6
		// Nothing happens

		currentmajor = 6
	}

	if migrated {
		logger.Info().WithFields(log.Fields{
			"from": config.FFVersion,
			"to":   "^" + version,
		}).Log("Migrated")
	}

	config.FFVersion = "^" + version

	return migrated, nil
}
