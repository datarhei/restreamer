package fs

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/datarhei/core/v16/io/fs"
	"github.com/datarhei/core/v16/log"
)

type Config struct {
	FS     fs.Filesystem
	Logger log.Logger
}

type Pattern struct {
	Pattern       string
	MaxFiles      uint
	MaxFileAge    time.Duration
	PurgeOnDelete bool
}

type Filesystem interface {
	fs.Filesystem

	// SetCleanup
	SetCleanup(id string, patterns []Pattern)

	// UnsetCleanup
	UnsetCleanup(id string)

	// Start
	Start()

	// Stop()
	Stop()
}

type filesystem struct {
	fs.Filesystem

	cleanupPatterns map[string][]Pattern
	cleanupLock     sync.RWMutex

	stopTicker context.CancelFunc

	startOnce sync.Once
	stopOnce  sync.Once

	logger log.Logger
}

func New(config Config) Filesystem {
	rfs := &filesystem{
		Filesystem: config.FS,
		logger:     config.Logger,
	}

	if rfs.logger == nil {
		rfs.logger = log.New("")
	}

	rfs.logger = rfs.logger.WithFields(log.Fields{
		"name": config.FS.Name(),
		"type": config.FS.Type(),
	})

	rfs.cleanupPatterns = make(map[string][]Pattern)

	// already drain the stop
	rfs.stopOnce.Do(func() {})

	return rfs
}

func (rfs *filesystem) Start() {
	rfs.startOnce.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		rfs.stopTicker = cancel
		go rfs.cleanupTicker(ctx, time.Second)

		rfs.stopOnce = sync.Once{}

		rfs.logger.Debug().Log("Starting cleanup")
	})
}

func (rfs *filesystem) Stop() {
	rfs.stopOnce.Do(func() {
		rfs.stopTicker()

		rfs.startOnce = sync.Once{}

		rfs.logger.Debug().Log("Stopping cleanup")
	})
}

func (rfs *filesystem) SetCleanup(id string, patterns []Pattern) {
	if len(patterns) == 0 {
		return
	}

	for _, p := range patterns {
		rfs.logger.Debug().WithFields(log.Fields{
			"id":           id,
			"pattern":      p.Pattern,
			"max_files":    p.MaxFiles,
			"max_file_age": p.MaxFileAge.Seconds(),
		}).Log("Add pattern")
	}

	rfs.cleanupLock.Lock()
	defer rfs.cleanupLock.Unlock()

	rfs.cleanupPatterns[id] = append(rfs.cleanupPatterns[id], patterns...)
}

func (rfs *filesystem) UnsetCleanup(id string) {
	rfs.logger.Debug().WithField("id", id).Log("Remove pattern group")

	rfs.cleanupLock.Lock()
	defer rfs.cleanupLock.Unlock()

	patterns := rfs.cleanupPatterns[id]
	delete(rfs.cleanupPatterns, id)

	rfs.purge(patterns)
}

func (rfs *filesystem) cleanup() {
	rfs.cleanupLock.RLock()
	defer rfs.cleanupLock.RUnlock()

	for _, patterns := range rfs.cleanupPatterns {
		for _, pattern := range patterns {
			filesAndDirs := rfs.Filesystem.List("/", pattern.Pattern)

			files := []fs.FileInfo{}
			for _, f := range filesAndDirs {
				if f.IsDir() {
					continue
				}

				files = append(files, f)
			}

			sort.Slice(files, func(i, j int) bool { return files[i].ModTime().Before(files[j].ModTime()) })

			if pattern.MaxFiles > 0 && uint(len(files)) > pattern.MaxFiles {
				for i := uint(0); i < uint(len(files))-pattern.MaxFiles; i++ {
					rfs.logger.Debug().WithField("path", files[i].Name()).Log("Remove file because MaxFiles is exceeded")
					rfs.Filesystem.Remove(files[i].Name())
				}
			}

			if pattern.MaxFileAge > 0 {
				bestBefore := time.Now().Add(-pattern.MaxFileAge)

				for _, f := range files {
					if f.ModTime().Before(bestBefore) {
						rfs.logger.Debug().WithField("path", f.Name()).Log("Remove file because MaxFileAge is exceeded")
						rfs.Filesystem.Remove(f.Name())
					}
				}
			}
		}
	}
}

func (rfs *filesystem) purge(patterns []Pattern) (nfiles uint64) {
	for _, pattern := range patterns {
		if !pattern.PurgeOnDelete {
			continue
		}

		files := rfs.Filesystem.List("/", pattern.Pattern)
		sort.Slice(files, func(i, j int) bool { return len(files[i].Name()) > len(files[j].Name()) })
		for _, f := range files {
			rfs.logger.Debug().WithField("path", f.Name()).Log("Purging file")
			rfs.Filesystem.Remove(f.Name())
			nfiles++
		}
	}

	return
}

func (rfs *filesystem) cleanupTicker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rfs.cleanup()
		}
	}
}
