package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/mansoormajeed/glimpse/internal/agent/grpcclient"
	"github.com/mansoormajeed/glimpse/internal/agent/heartbeat"
	"github.com/mansoormajeed/glimpse/internal/common/logger"
)

func main() {

	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	if *debug {
		logger.SetDebugLevel()
		logger.Debug("Debug logging enabled")
	}

	hostname, err := os.Hostname()
	if err != nil {
		logger.Errorf("Error getting hostname: %v", err)
		return
	}
	logger.Infof("Starting the agent on host: %s", hostname)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupSignalHandling(cancel)
	run(ctx)

	logger.Infof("Agent is running... Press Ctrl+C to exit.")

	// Wait for the context to be cancelled
	<-ctx.Done()
}

func run(ctx context.Context) {

	grpcClient, err := grpcclient.NewGlimpseServiceClient()
	if err != nil {
		logger.Errorf("Error creating gRPC client: %v", err)
		return
	}
	heartbeatService := heartbeat.NewHeartbeatService(grpcClient)
}

func setupSignalHandling(cancel context.CancelFunc) {
	// handle OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logger.Infof("Received signal: %v, shutting down..", sig)
		cancel()
	}()
}
