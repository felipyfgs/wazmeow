package sqlite

import (
	"github.com/uptrace/bun"

	"wazmeow/internal/infra/config"
	"wazmeow/internal/infra/database"
	"wazmeow/pkg/logger"
)

// Connection represents a SQLite database connection with repositories
type Connection struct {
	db          database.Connection
	sessionRepo *SessionRepository
	migrator    interface{} // Will be replaced with actual migrator
	logger      logger.Logger
}

// New creates a new SQLite connection with all repositories
func New(cfg *config.DatabaseConfig, log logger.Logger) (*Connection, error) {
	// Create database connection
	dbConn, err := database.NewSQLite(cfg, log)
	if err != nil {
		return nil, err
	}

	// Create migrator (placeholder for now)
	var migrator interface{} = nil

	// Auto-migrate if enabled (placeholder for now)
	if cfg.AutoMigrate {
		log.Info("auto-migration is enabled but not implemented yet")
		// TODO: Implement migration logic
	}

	// Create repositories
	bunDB := dbConn.GetDB()
	sessionRepo := NewSessionRepository(bunDB, log).(*SessionRepository)

	return &Connection{
		db:          dbConn,
		sessionRepo: sessionRepo,
		migrator:    migrator,
		logger:      log,
	}, nil
}

// GetDB returns the underlying Bun DB instance
func (c *Connection) GetDB() *bun.DB {
	return c.db.GetDB()
}

// SessionRepository returns the session repository
func (c *Connection) SessionRepository() *SessionRepository {
	return c.sessionRepo
}

// Migrator returns the database migrator (placeholder)
func (c *Connection) Migrator() interface{} {
	return c.migrator
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}

// Health checks the database health
func (c *Connection) Health() error {
	return c.db.Health()
}
