package app

import (
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Data      string
}

type LogHistoryEntry struct {
	CreatedAt time.Time
	Prelude   []string
	Log       []LogEntry
}

type Log struct {
	LogHistoryEntry
	History []LogHistoryEntry
}
