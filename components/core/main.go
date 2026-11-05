package main

import (
	"os"
	"os/signal"

	"github.com/datarhei/core/v16/app/api"
	"github.com/datarhei/core/v16/config/store"
	"github.com/datarhei/core/v16/log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := log.New("Core").WithOutput(log.NewConsoleWriter(os.Stderr, log.Lwarn, true))

	configfile := store.Location(os.Getenv("CORE_CONFIGFILE"))

	app, err := api.New(configfile, os.Stderr)
	if err != nil {
		logger.Error().WithError(err).Log("Failed to create new API")
		os.Exit(1)
	}

	go func() {
		defer func() {
			if proc, err := os.FindProcess(os.Getpid()); err == nil {
				proc.Signal(os.Interrupt)
			}
		}()

		for {
			if err := app.Start(); err != api.ErrConfigReload {
				if err != nil {
					logger.Error().WithError(err).Log("Failed to start API")
				}

				break
			} else {
				logger.Warn().WithError(err).Log("Config reload requested")
			}

			app.Stop()

			if err := app.Reload(); err != nil {
				logger.Error().WithError(err).Log("Failed to reload config")
				break
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the app
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Stop the app
	app.Destroy()
}
