package main

import (
	"os/signal"
	"time"
)

func main() {
	// Ignore all signals
	signal.Ignore()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
	}
}
