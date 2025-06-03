package main

import (
	"backend/internal/aggregator"
	"backend/internal/env_var"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	slog.Info("Ensuring that environment are set")
	e := env_var.GetAggVars()

	slog.Info("Starting server...", "port", e.AggPort)

	// Setup listener for signal interrupts etc.
	c := make(chan os.Signal, 1)
	shutdownWg := sync.WaitGroup{}
	shutdownWg.Add(1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		slog.Warn("Received OS signal, stopping server", "signal", sig)
		// Perform cleanup
		aggregator.GetAggregator().CleanUp()
		shutdownWg.Done()
	}()

	// Start server
	aggregator.Serve(e.AggPort)
	shutdownWg.Wait()

}
