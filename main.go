package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/my-desktop-agent/agent" // Your agent logic
	"github.com/yourusername/my-desktop-agent/config"
	"github.com/yourusername/my-desktop-agent/logger"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load configuration
	cfg, err := config.LoadConfig("config.json") // Or from environment variables, etc.
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create agent instance
	myAgent := agent.NewAgent(cfg)

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown (Ctrl+C, etc.)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Received shutdown signal, shutting down agent...")
		cancel()
	}()

	// Start the agent's main loop
	if err := myAgent.Start(ctx); err != nil {
		log.Fatalf("Agent failed to start: %v", err)
	}

	// Keep main goroutine alive until shutdown
	<-ctx.Done()
	log.Println("Agent stopped.")
}
