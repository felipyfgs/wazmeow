package app

import (
	"context"
	"fmt"
	"time"

	"wazmeow/internal/http/handler"
	"wazmeow/internal/http/routes"
	"wazmeow/internal/http/server"
	"wazmeow/internal/infra/config"
	"wazmeow/internal/infra/container"
	sessionUC "wazmeow/internal/usecases/session"
	whatsappUC "wazmeow/internal/usecases/whatsapp"
	"wazmeow/pkg/logger"
)

// AppContainer holds all application dependencies
type AppContainer struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	InfraContainer *container.Container

	// Use Cases
	SessionCreateUC        *sessionUC.CreateUseCase
	SessionConnectUC       *sessionUC.ConnectUseCase
	SessionDisconnectUC    *sessionUC.DisconnectUseCase
	SessionListUC          *sessionUC.ListUseCase
	SessionDeleteUC        *sessionUC.DeleteUseCase
	SessionResolveUC       *sessionUC.ResolveUseCase
	SessionAutoReconnectUC *sessionUC.AutoReconnectUseCase

	WhatsAppGenerateQRUC  *whatsappUC.GenerateQRUseCase
	WhatsAppPairPhoneUC   *whatsappUC.PairPhoneUseCase
	WhatsAppSendMessageUC *whatsappUC.SendMessageUseCase

	// HTTP Layer
	SessionHandler *handler.SessionHandler
	HealthHandler  *handler.HealthHandler
	Router         *routes.Router
	HTTPServer     *server.Server
	ServerManager  *server.ServerManager

	// Internal state
	isInitialized bool
}

// NewAppContainer creates a new application container
func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	container := &AppContainer{
		Config: cfg,
	}

	if err := container.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize app container: %w", err)
	}

	return container, nil
}

// initialize sets up all application dependencies
func (c *AppContainer) initialize() error {
	// Initialize infrastructure container
	if err := c.initializeInfrastructure(); err != nil {
		return fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	// Initialize use cases
	if err := c.initializeUseCases(); err != nil {
		return fmt.Errorf("failed to initialize use cases: %w", err)
	}

	// Initialize HTTP layer
	if err := c.initializeHTTPLayer(); err != nil {
		return fmt.Errorf("failed to initialize HTTP layer: %w", err)
	}

	c.isInitialized = true
	c.InfraContainer.Logger.Info("Application container initialized successfully")

	return nil
}

// initializeInfrastructure sets up infrastructure dependencies
func (c *AppContainer) initializeInfrastructure() error {
	infraContainer, err := container.New(c.Config)
	if err != nil {
		return fmt.Errorf("failed to create infrastructure container: %w", err)
	}

	c.InfraContainer = infraContainer
	return nil
}

// initializeUseCases sets up use case dependencies
func (c *AppContainer) initializeUseCases() error {
	logger := c.InfraContainer.Logger
	validator := c.InfraContainer.Validator

	// Session use cases
	c.SessionCreateUC = sessionUC.NewCreateUseCase(
		c.InfraContainer.SessionRepo,
		logger,
		validator,
	)

	c.SessionConnectUC = sessionUC.NewConnectUseCase(
		c.InfraContainer.SessionRepo,
		c.InfraContainer.WhatsAppManager,
		logger,
	)

	c.SessionDisconnectUC = sessionUC.NewDisconnectUseCase(
		c.InfraContainer.SessionRepo,
		c.InfraContainer.WhatsAppManager,
		logger,
	)

	c.SessionListUC = sessionUC.NewListUseCase(
		c.InfraContainer.SessionRepo,
		logger,
	)

	c.SessionDeleteUC = sessionUC.NewDeleteUseCase(
		c.InfraContainer.SessionRepo,
		c.InfraContainer.WhatsAppManager,
		logger,
	)

	c.SessionResolveUC = sessionUC.NewResolveUseCase(
		c.InfraContainer.SessionRepo,
		logger,
	)

	c.SessionAutoReconnectUC = sessionUC.NewAutoReconnectUseCase(
		c.InfraContainer.SessionRepo,
		c.InfraContainer.WhatsAppManager,
		logger,
	)

	// WhatsApp use cases
	c.WhatsAppGenerateQRUC = whatsappUC.NewGenerateQRUseCase(
		c.InfraContainer.SessionRepo,
		c.InfraContainer.WhatsAppManager,
		logger,
	)

	c.WhatsAppPairPhoneUC = whatsappUC.NewPairPhoneUseCase(
		c.InfraContainer.SessionRepo,
		c.InfraContainer.WhatsAppManager,
		logger,
		validator,
	)

	c.WhatsAppSendMessageUC = whatsappUC.NewSendMessageUseCase(
		c.InfraContainer.SessionRepo,
		c.InfraContainer.WhatsAppManager,
		logger,
		validator,
	)

	logger.Info("Use cases initialized")
	return nil
}

// initializeHTTPLayer sets up HTTP layer dependencies
func (c *AppContainer) initializeHTTPLayer() error {
	logger := c.InfraContainer.Logger
	validator := c.InfraContainer.Validator

	// Create handlers
	c.SessionHandler = handler.NewSessionHandler(
		c.SessionCreateUC,
		c.SessionConnectUC,
		c.SessionDisconnectUC,
		c.SessionListUC,
		c.SessionDeleteUC,
		c.SessionResolveUC,
		c.WhatsAppGenerateQRUC,
		c.WhatsAppPairPhoneUC,
		logger,
		validator,
	)

	c.HealthHandler = handler.NewHealthHandler(
		c.InfraContainer,
		logger,
	)

	// Create router
	c.Router = routes.NewRouter(
		c.SessionHandler,
		c.HealthHandler,
		c.Config,
		logger,
	)

	// Create HTTP server
	c.HTTPServer = server.New(
		c.Router,
		&c.Config.Server,
		logger,
	)

	// Create server manager
	c.ServerManager = server.NewServerManager(
		c.HTTPServer,
		logger,
	)

	logger.Info("HTTP layer initialized")
	return nil
}

// GetLogger returns the application logger
func (c *AppContainer) GetLogger() logger.Logger {
	if c.InfraContainer != nil {
		return c.InfraContainer.Logger
	}
	return nil
}

// GetConfig returns the application configuration
func (c *AppContainer) GetConfig() *config.Config {
	return c.Config
}

// GetServerManager returns the server manager
func (c *AppContainer) GetServerManager() *server.ServerManager {
	return c.ServerManager
}

// GetInfraContainer returns the infrastructure container
func (c *AppContainer) GetInfraContainer() *container.Container {
	return c.InfraContainer
}

// IsInitialized returns true if the container is initialized
func (c *AppContainer) IsInitialized() bool {
	return c.isInitialized
}

// Close gracefully shuts down all dependencies
func (c *AppContainer) Close() error {
	if !c.isInitialized {
		return nil
	}

	c.InfraContainer.Logger.Info("Shutting down application container")

	var errors []error

	// Close infrastructure container
	if c.InfraContainer != nil {
		if err := c.InfraContainer.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close infrastructure container: %w", err))
		}
	}

	if len(errors) > 0 {
		// Log all errors
		for _, err := range errors {
			c.InfraContainer.Logger.ErrorWithError("error during app container shutdown", err, nil)
		}
		return fmt.Errorf("multiple errors during shutdown: %v", errors)
	}

	c.InfraContainer.Logger.Info("Application container shut down successfully")
	return nil
}

// Health checks the health of all application components
func (c *AppContainer) Health() error {
	if !c.isInitialized {
		return fmt.Errorf("container not initialized")
	}

	// Check infrastructure health
	if err := c.InfraContainer.Health(); err != nil {
		return fmt.Errorf("infrastructure health check failed: %w", err)
	}

	// Check HTTP server health
	if err := c.HTTPServer.Health(); err != nil {
		return fmt.Errorf("HTTP server health check failed: %w", err)
	}

	return nil
}

// GetServerInfo returns information about the HTTP server
func (c *AppContainer) GetServerInfo() server.ServerInfo {
	if c.ServerManager != nil {
		return c.ServerManager.GetServerInfo()
	}
	return server.ServerInfo{}
}

// StartWhatsAppManager starts the WhatsApp manager
func (c *AppContainer) StartWhatsAppManager() error {
	return c.InfraContainer.StartWhatsAppManager()
}

// GetDatabaseStats returns database statistics
func (c *AppContainer) GetDatabaseStats() interface{} {
	return c.InfraContainer.GetDatabaseStats()
}

// GetWhatsAppStats returns WhatsApp statistics
func (c *AppContainer) GetWhatsAppStats() interface{} {
	return c.InfraContainer.GetWhatsAppStats()
}

// AutoReconnectSessions performs automatic reconnection of eligible sessions
func (c *AppContainer) AutoReconnectSessions(ctx context.Context) (*sessionUC.AutoReconnectResponse, error) {
	if !c.isInitialized {
		return nil, fmt.Errorf("container not initialized")
	}

	req := sessionUC.AutoReconnectRequest{
		MaxConcurrentReconnections: 5, // Limit concurrent reconnections
		ReconnectionTimeout:        30 * time.Second,
	}

	return c.SessionAutoReconnectUC.Execute(ctx, req)
}
