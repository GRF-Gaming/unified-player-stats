package main

import (
	"backend/internal/aggregator"
	"backend/internal/utils/env_var"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	slog.Info("Ensuring that aggregator environment variables are set")
	e := env_var.GetAggVars()

	// Setting aggregator service
	slog.Info("Starting aggregator...", "port", e.AggPort)
	a := aggregator.GetAggregator()

	// Setup listener for signal interrupts etc.
	c := make(chan os.Signal, 1)
	shutdownWg := sync.WaitGroup{}
	shutdownWg.Add(1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		slog.Warn("Received OS signal, stopping server", "signal", sig)
		// Perform cleanup
		a.CleanUp()
		shutdownWg.Done()
	}()

	// Start server
	a.Serve(e.AggPort)
	shutdownWg.Wait()

}
