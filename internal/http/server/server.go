package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"wazmeow/internal/http/routes"
	"wazmeow/internal/infra/config"
	"wazmeow/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	httpServer *http.Server
	router     *routes.Router
	config     *config.ServerConfig
	logger     logger.Logger
}

// New creates a new HTTP server
func New(router *routes.Router, config *config.ServerConfig, logger logger.Logger) *Server {
	return &Server{
		router: router,
		config: config,
		logger: logger,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Setup routes
	handler := s.router.SetupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      handler,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	s.logger.InfoWithFields("Starting HTTP server", logger.Fields{
		"host":          s.config.Host,
		"port":          s.config.Port,
		"read_timeout":  s.config.ReadTimeout,
		"write_timeout": s.config.WriteTimeout,
		"idle_timeout":  s.config.IdleTimeout,
	})

	// Start server
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Info("Stopping HTTP server...")

	// Attempt graceful shutdown
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.ErrorWithError("Failed to gracefully shutdown HTTP server", err, nil)

		// Force close if graceful shutdown fails
		if closeErr := s.httpServer.Close(); closeErr != nil {
			return fmt.Errorf("failed to force close HTTP server: %w", closeErr)
		}

		return fmt.Errorf("failed to gracefully shutdown HTTP server: %w", err)
	}

	s.logger.Info("HTTP server stopped gracefully")
	return nil
}

// GetAddr returns the server address
func (s *Server) GetAddr() string {
	if s.httpServer != nil {
		return s.httpServer.Addr
	}
	return fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
}

// IsRunning returns true if the server is running
func (s *Server) IsRunning() bool {
	return s.httpServer != nil
}


// Health checks if the server is healthy
func (s *Server) Health() error {
	if !s.IsRunning() {
		return fmt.Errorf("server not running")
	}
	return nil
}

// ServerManager manages the HTTP server lifecycle
type ServerManager struct {
	server *Server
	logger logger.Logger
}

// NewServerManager creates a new server manager
func NewServerManager(server *Server, logger logger.Logger) *ServerManager {
	return &ServerManager{
		server: server,
		logger: logger,
	}
}

// StartWithGracefulShutdown starts the server and handles graceful shutdown
func (sm *ServerManager) StartWithGracefulShutdown(ctx context.Context) error {
	// Channel to receive server errors
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		sm.logger.InfoWithFields("HTTP server starting", logger.Fields{
			"addr": sm.server.GetAddr(),
		})

		if err := sm.server.Start(); err != nil {
			serverErrors <- err
		}
	}()

	// Wait for context cancellation or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		sm.logger.Info("Shutdown signal received, stopping HTTP server...")

		// Create shutdown context with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Stop server gracefully
		if err := sm.server.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}

		sm.logger.Info("HTTP server shutdown completed")
		return nil
	}
}

// GetServerInfo returns information about the server
func (sm *ServerManager) GetServerInfo() ServerInfo {
	return ServerInfo{
		Address:   sm.server.GetAddr(),
		IsRunning: sm.server.IsRunning(),
	}
}

// ServerInfo holds information about the server
type ServerInfo struct {
	Address   string `json:"address"`
	IsRunning bool   `json:"is_running"`
}

// HealthCheckServer provides health check functionality
type HealthCheckServer struct {
	server *http.Server
	logger logger.Logger
}

// NewHealthCheckServer creates a dedicated health check server
func NewHealthCheckServer(port int, logger logger.Logger) *HealthCheckServer {
	mux := http.NewServeMux()

	// Simple health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &HealthCheckServer{
		server: server,
		logger: logger,
	}
}

// Start starts the health check server
func (hcs *HealthCheckServer) Start() error {
	hcs.logger.InfoWithFields("Starting health check server", logger.Fields{
		"addr": hcs.server.Addr,
	})

	if err := hcs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start health check server: %w", err)
	}

	return nil
}

// Stop stops the health check server
func (hcs *HealthCheckServer) Stop(ctx context.Context) error {
	hcs.logger.Info("Stopping health check server...")

	if err := hcs.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop health check server: %w", err)
	}

	hcs.logger.Info("Health check server stopped")
	return nil
}

// StartHealthCheckServer starts a health check server in a separate goroutine
func StartHealthCheckServer(port int, logger logger.Logger) *HealthCheckServer {
	healthServer := NewHealthCheckServer(port, logger)

	go func() {
		if err := healthServer.Start(); err != nil {
			fields := map[string]interface{}{
				"port": port,
			}
			logger.ErrorWithError("Health check server error", err, fields)
		}
	}()

	return healthServer
}
