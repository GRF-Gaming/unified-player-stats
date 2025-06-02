package main

import (
	"backend/internal/aggregator"
	"backend/internal/models"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.Info("Starting server...")

	key, _ := models.NewApiKeyComponentsV1().ToApiKeyV1()
	slog.Info("TESTING", "key", key)

	// Setup listener for signal interrupts etc.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		slog.Warn("Received OS signal, stopping server", "signal", sig)
		// Perform cleanup
	}()

	// Start server
	aggregator.Serve(8080)

}
