package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zinrai/sevalet/internal/config"
)

type Server struct {
	config *config.Config
	server *http.Server
}

func NewServer(config *config.Config) *Server {
	handler := SetupRoutes(config)

	return &Server{
		config: config,
		server: &http.Server{
			Addr:    config.GetAddress(),
			Handler: handler,
		},
	}
}

// Start launches the server
func (s *Server) Start() error {
	// Create signal reception channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Create error channel
	errChan := make(chan error, 1)

	// Start server in a separate goroutine
	go func() {
		log.Printf("Starting server: %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Wait for signal or error
	select {
	case <-stop:
		log.Println("Received shutdown signal")
	case err := <-errChan:
		return fmt.Errorf("failed to start server: %w", err)
	}

	log.Println("Shutting down server...")

	// Shutdown with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	log.Println("Server shutdown completed")
	return nil
}
