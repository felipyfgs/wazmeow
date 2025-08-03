package database

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"

	"wazmeow/internal/infra/config"
	"wazmeow/internal/infra/database/drivers"
	"wazmeow/pkg/logger"
)

// DriverConnection interface for driver implementations
type DriverConnection interface {
	GetDB() *bun.DB
	Close() error
	Health() error
	Stats() sql.DBStats
}

// connectionAdapter adapts driver connections to database.Connection interface
type connectionAdapter struct {
	conn DriverConnection
}

func (a *connectionAdapter) GetDB() *bun.DB {
	return a.conn.GetDB()
}

func (a *connectionAdapter) Close() error {
	return a.conn.Close()
}

func (a *connectionAdapter) Health() error {
	return a.conn.Health()
}

func (a *connectionAdapter) Stats() sql.DBStats {
	return a.conn.Stats()
}

// ConnectionFactory creates database connections based on configuration
type ConnectionFactory struct {
	logger logger.Logger
}

// NewConnectionFactory creates a new connection factory
func NewConnectionFactory(log logger.Logger) *ConnectionFactory {
	return &ConnectionFactory{
		logger: log,
	}
}

// CreateConnection creates a database connection based on the driver type
func (f *ConnectionFactory) CreateConnection(cfg *config.DatabaseConfig) (Connection, error) {
	dbType := DatabaseType(cfg.Driver)

	f.logger.InfoWithFields("creating database connection", logger.Fields{
		"driver": cfg.Driver,
		"type":   string(dbType),
	})

	switch dbType {
	case SQLite, SQLite3:
		conn, err := drivers.NewSQLiteConnection(cfg, f.logger)
		if err != nil {
			return nil, err
		}
		return &connectionAdapter{conn: conn}, nil
	case PostgreSQL, PostgreSQL2:
		conn, err := drivers.NewPostgreSQLConnection(cfg, f.logger)
		if err != nil {
			return nil, err
		}
		return &connectionAdapter{conn: conn}, nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

// New creates a new database connection using the factory pattern
func New(cfg *config.DatabaseConfig, log logger.Logger) (Connection, error) {
	factory := NewConnectionFactory(log)
	return factory.CreateConnection(cfg)
}

// NewSQLite creates a new SQLite database connection
func NewSQLite(cfg *config.DatabaseConfig, log logger.Logger) (Connection, error) {
	conn, err := drivers.NewSQLiteConnection(cfg, log)
	if err != nil {
		return nil, err
	}
	return &connectionAdapter{conn: conn}, nil
}

// NewPostgreSQL creates a new PostgreSQL database connection
func NewPostgreSQL(cfg *config.DatabaseConfig, log logger.Logger) (Connection, error) {
	conn, err := drivers.NewPostgreSQLConnection(cfg, log)
	if err != nil {
		return nil, err
	}
	return &connectionAdapter{conn: conn}, nil
}
