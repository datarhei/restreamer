package store

import (
	gojson "encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/datarhei/core/v16/encoding/json"
	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/log"
)

type JSONConfig struct {
	Filesystem fs.Filesystem
	Filepath   string // Full path to the database file
	Logger     log.Logger
}

type jsonStore struct {
	fs       fs.Filesystem
	filepath string
	logger   log.Logger

	// Mutex to serialize access to the backend
	lock sync.RWMutex
}

var version uint64 = 4

func NewJSON(config JSONConfig) (Store, error) {
	s := &jsonStore{
		fs:       config.Filesystem,
		filepath: config.Filepath,
		logger:   config.Logger,
	}

	if len(s.filepath) == 0 {
		s.filepath = "/db.json"
	}

	if s.fs == nil {
		return nil, fmt.Errorf("no valid filesystem provided")
	}

	if s.logger == nil {
		s.logger = log.New("")
	}

	return s, nil
}

func (s *jsonStore) Load() (StoreData, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	data, err := s.load(s.filepath, version)
	if err != nil {
		return NewStoreData(), err
	}

	data.sanitize()

	return data, nil
}

func (s *jsonStore) Store(data StoreData) error {
	if data.Version != version {
		return fmt.Errorf("invalid version (have: %d, want: %d)", data.Version, version)
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	err := s.store(s.filepath, data)
	if err != nil {
		return fmt.Errorf("failed to store data: %w", err)
	}

	return nil
}

func (s *jsonStore) store(filepath string, data StoreData) error {
	jsondata, err := gojson.MarshalIndent(&data, "", "    ")
	if err != nil {
		return err
	}

	_, _, err = s.fs.WriteFileSafe(filepath, jsondata)
	if err != nil {
		return err
	}

	s.logger.WithField("file", filepath).Debug().Log("Stored data")

	return nil
}

type storeVersion struct {
	Version uint64 `json:"version"`
}

func (s *jsonStore) load(filepath string, version uint64) (StoreData, error) {
	r := NewStoreData()

	_, err := s.fs.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return r, nil
		}

		return r, err
	}

	jsondata, err := s.fs.ReadFile(filepath)
	if err != nil {
		return r, err
	}

	var db storeVersion

	if err = gojson.Unmarshal(jsondata, &db); err != nil {
		return r, json.FormatError(jsondata, err)
	}

	if db.Version != version {
		return r, fmt.Errorf("unsupported version of the DB file (want: %d, have: %d)", version, db.Version)
	}

	if err = gojson.Unmarshal(jsondata, &r); err != nil {
		return r, json.FormatError(jsondata, err)
	}

	s.logger.WithField("file", filepath).Debug().Log("Read data")

	return r, nil
}
