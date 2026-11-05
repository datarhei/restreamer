package store

import "github.com/datarhei/core/v16/config"

// Store is a store for the configuration data.
type Store interface {
	// Get the current configuration.
	Get() *config.Config

	// Set a new configuration for persistence.
	Set(data *config.Config) error

	// GetActive returns the configuration that has been set as
	// active before, otherwise it return nil.
	GetActive() *config.Config

	// SetActive will keep the given configuration
	// as active in memory. It can be retrieved later with GetActive()
	SetActive(data *config.Config) error

	// Reload will reload the stored configuration. It has to make sure
	// that all affected components will receiver their potentially
	// changed configuration.
	Reload() error
}

type DataVersion struct {
	Version int64 `json:"version"`
}
