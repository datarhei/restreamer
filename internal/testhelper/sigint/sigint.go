package main

import (
	"os"
	"os/signal"
)

func main() {
	// Wait for interrupt signal to gracefully shutdown the app
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	os.Exit(255)
}
