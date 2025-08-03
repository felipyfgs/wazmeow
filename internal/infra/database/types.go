package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	"wazmeow/internal/infra/config"
	"wazmeow/pkg/logger"
)

// Driver interface defines database driver implementations
type Driver interface {
	Connect(cfg *config.DatabaseConfig, log logger.Logger) (Connection, error)
	Name() string
	SupportedDrivers() []string
}

// Migrator interface defines database migration operations
type Migrator interface {
	Migrate(ctx context.Context) error
	Reset(ctx context.Context) error
	Drop(ctx context.Context) error
}

// DatabaseType represents the type of database
type DatabaseType string

const (
	SQLite      DatabaseType = "sqlite"
	SQLite3     DatabaseType = "sqlite3"
	PostgreSQL  DatabaseType = "postgres"
	PostgreSQL2 DatabaseType = "postgresql"
)

// BaseConnection represents the base database connection structure
type BaseConnection struct {
	DB     *bun.DB
	Config *config.DatabaseConfig
	Logger logger.Logger
}

// Close closes the database connection
func (c *BaseConnection) Close() error {
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			c.Logger.ErrorWithError("failed to close database connection", err, nil)
			return err
		}
		c.Logger.Info("database connection closed")
	}
	return nil
}

// GetDB returns the Bun DB instance
func (c *BaseConnection) GetDB() *bun.DB {
	return c.DB
}

// Health checks the database health
func (c *BaseConnection) Health() error {
	if c.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := c.DB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// Stats returns database connection statistics
func (c *BaseConnection) Stats() sql.DBStats {
	if c.DB == nil {
		return sql.DBStats{}
	}
	return c.DB.DB.Stats()
}
