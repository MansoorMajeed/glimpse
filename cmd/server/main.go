package main

import (
	"flag"

	"github.com/mansoormajeed/glimpse/internal/logger"
)

func main() {

	listenPort := flag.Int("port", 5001, "Port to listen on")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetDebugLevel()
		logger.Debug("Debug logging enabled")
	}
	logger.Info("Starting the server...")
	logger.Debugf("Listening on port: %d", *listenPort)

}
