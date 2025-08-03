package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wazmeow/internal/infra/config"
	"wazmeow/pkg/logger"
)

// App represents the main application
type App struct {
	container *AppContainer
	logger    logger.Logger
}

// New creates a new application instance
func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create application container
	appContainer, err := NewAppContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create app container: %w", err)
	}

	return &App{
		container: appContainer,
		logger:    appContainer.GetLogger(),
	}, nil
}

// Start starts the application with graceful shutdown
func (a *App) Start() error {
	a.logger.Info("Starting WazMeow application")

	// Start WhatsApp manager
	if err := a.container.StartWhatsAppManager(); err != nil {
		return fmt.Errorf("failed to start WhatsApp manager: %w", err)
	}

	// Perform automatic reconnection of sessions
	a.logger.Info("🔄 Starting automatic session reconnection")
	reconnectCtx, reconnectCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer reconnectCancel()

	reconnectResult, err := a.container.AutoReconnectSessions(reconnectCtx)
	if err != nil {
		a.logger.ErrorWithError("failed to perform automatic reconnection", err, nil)
		// Don't fail the startup, just log the error
	} else {
		a.logger.InfoWithFields("automatic reconnection completed", logger.Fields{
			"total_sessions":           reconnectResult.TotalSessions,
			"eligible_sessions":        reconnectResult.EligibleSessions,
			"successful_reconnections": reconnectResult.SuccessfulReconnections,
			"failed_reconnections":     reconnectResult.FailedReconnections,
		})

		// Log individual session results
		for _, result := range reconnectResult.ReconnectionResults {
			if result.Success {
				a.logger.InfoWithFields("✅ session reconnected", logger.Fields{
					"session_id":   result.SessionID.String(),
					"session_name": result.SessionName,
					"duration_ms":  result.Duration.Milliseconds(),
				})
			} else {
				a.logger.WarnWithFields("❌ session reconnection failed", logger.Fields{
					"session_id":   result.SessionID.String(),
					"session_name": result.SessionName,
					"error":        result.Error,
					"duration_ms":  result.Duration.Milliseconds(),
				})
			}
		}
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start HTTP server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		a.logger.InfoWithFields("Starting HTTP server", logger.Fields{
			"address": a.container.HTTPServer.GetAddr(),
		})

		if err := a.container.ServerManager.StartWithGracefulShutdown(ctx); err != nil {
			serverErrors <- err
		}
	}()

	a.logger.InfoWithFields("WazMeow application started successfully", logger.Fields{
		"server_address": a.container.HTTPServer.GetAddr(),
	})

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		a.logger.ErrorWithError("Server error", err, nil)
		return fmt.Errorf("server error: %w", err)
	case sig := <-sigChan:
		a.logger.InfoWithFields("Shutdown signal received", logger.Fields{
			"signal": sig.String(),
		})
		cancel() // Cancel context to trigger graceful shutdown
		return nil
	}
}

// Stop stops the application
func (a *App) Stop() error {
	a.logger.Info("Stopping WazMeow application")

	if err := a.container.Close(); err != nil {
		return fmt.Errorf("failed to close app container: %w", err)
	}

	a.logger.Info("WazMeow application stopped successfully")
	return nil
}

// Health checks the application health
func (a *App) Health() error {
	return a.container.Health()
}

// GetConfig returns the application configuration
func (a *App) GetConfig() *config.Config {
	return a.container.GetConfig()
}

// GetContainer returns the application container
func (a *App) GetContainer() *AppContainer {
	return a.container
}

// GetServerInfo returns information about the HTTP server
func (a *App) GetServerInfo() interface{} {
	return a.container.GetServerInfo()
}

// GetStats returns application statistics
func (a *App) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"server":   a.container.GetServerInfo(),
		"database": a.container.GetDatabaseStats(),
		"whatsapp": a.container.GetWhatsAppStats(),
	}
}
