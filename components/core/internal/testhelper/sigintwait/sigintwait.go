package main

import (
	"os"
	"os/signal"
	"time"
)

func main() {
	// Wait for interrupt signal to gracefully shutdown the app
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	time.Sleep(3 * time.Second)

	os.Exit(255)
}
