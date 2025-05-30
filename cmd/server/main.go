package main

import (
	"flag"

	"github.com/mansoormajeed/glimpse/internal/common/logger"
	"github.com/mansoormajeed/glimpse/internal/server"
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

	store := server.NewServerStore(5) // for now use buffer 5. later change to 60

	server.StartHTTPServer(store)

	// Start the gRPC server
	server.StartGRPCServer(store)

}
