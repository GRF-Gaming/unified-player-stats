package main

import (
	"backend/internal/env_var"
	"backend/internal/processor_kills"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	slog.Info("Ensuring processor environment variable are set")
	_ = env_var.GetProVars()

	slog.Info("Starting processor...")
	p := processor_kills.GetProcessor()

	// Setup listener for signal interrupts etc.
	c := make(chan os.Signal, 1)
	shutdownWg := sync.WaitGroup{}
	shutdownWg.Add(1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		slog.Warn("Received OS signal, stopping server", "signal", sig)
		// Perform cleanup
		p.CleanUp()
		shutdownWg.Done()
	}()

	// Start processor
	p.Spin()
	shutdownWg.Wait()
}
