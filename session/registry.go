package session

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/log"
)

// Config is the configuration for creating a new registry
type Config struct {
	// PersistFS is a filesystem in whose root the session history will be persisted. If it is nil, the
	// history will not be persisted.
	PersistFS fs.Filesystem

	// Logger is an instance of a logger. If it is nil, no logs
	// will be written.
	Logger log.Logger
}

type RegistryReader interface {
	// Collectors returns an array of all registered IDs
	Collectors() []string

	// Collector returns the collector with the ID, or nil if the ID is not registered
	Collector(id string) Collector

	// Summary returns the summary from a collector with the ID, or an empty summary if the ID is not registered
	Summary(id string) Summary

	// Active returns the active sessions from a collector with the ID, or an empty list if the ID is not registered
	Active(id string) []Session
}

// The Registry interface
type Registry interface {
	// Register returns a new collector from conf and registers it under the id and error is nil. In case of error
	// the returned collector is nil and the error is not nil.
	Register(id string, conf CollectorConfig) (Collector, error)

	// Unregister unregisters the collector with the ID, returns error if the ID is not registered
	Unregister(id string) error

	// UnregisterAll unregisters al registered collectors
	UnregisterAll()

	RegistryReader
}

type registry struct {
	collector map[string]*collector
	persistFS fs.Filesystem
	logger    log.Logger

	lock sync.Mutex
}

// New returns a new registry for collectors that implement the Registry interface. The error
// is non-nil if the PersistDir from the config can't be created.
func New(conf Config) (Registry, error) {
	r := &registry{
		collector: make(map[string]*collector),
		persistFS: conf.PersistFS,
		logger:    conf.Logger,
	}

	if r.logger == nil {
		r.logger = log.New("Session")
	}

	return r, nil
}

func (r *registry) Register(id string, conf CollectorConfig) (Collector, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("invalid ID. empty IDs are not allowed")
	}

	re := regexp.MustCompile(`^[0-9A-Za-z]+$`)
	if !re.MatchString(id) {
		return nil, fmt.Errorf("invalid ID. only numbers and letters are allowed")
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.collector[id]
	if ok {
		return nil, fmt.Errorf("a collector with the ID '%s' already exists", id)
	}

	m, err := newCollector(id, r.persistFS, r.logger, conf)
	if err != nil {
		return nil, err
	}

	m.start()

	r.collector[id] = m

	return m, nil
}

func (r *registry) Unregister(id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	m, ok := r.collector[id]
	if !ok {
		return fmt.Errorf("a collector with the ID '%s' doesn't exist", id)
	}

	m.Stop()

	delete(r.collector, id)

	return nil
}

func (r *registry) Collectors() []string {
	collectors := []string{}

	r.lock.Lock()
	defer r.lock.Unlock()

	for name := range r.collector {
		collectors = append(collectors, name)
	}

	return collectors
}

func (r *registry) Collector(id string) Collector {
	r.lock.Lock()
	defer r.lock.Unlock()

	m, ok := r.collector[id]
	if !ok {
		return nil
	}

	return m
}

func (r *registry) UnregisterAll() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, m := range r.collector {
		m.Stop()
	}

	r.collector = make(map[string]*collector)
}

func (r *registry) Summary(id string) Summary {
	r.lock.Lock()
	defer r.lock.Unlock()

	m, ok := r.collector[id]
	if !ok {
		return Summary{}
	}

	return m.Summary()
}

func (r *registry) Active(id string) []Session {
	r.lock.Lock()
	defer r.lock.Unlock()

	m, ok := r.collector[id]
	if !ok {
		return []Session{}
	}

	return m.Active()
}
