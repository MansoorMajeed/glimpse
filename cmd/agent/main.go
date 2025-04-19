package main

import (
	"flag"

	"github.com/mansoormajeed/glimpse/internal/logger"
)

func main() {

	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetDebugLevel()
		logger.Debug("Debug logging enabled")
	}

	logger.Info("Starting the agent...")
}
