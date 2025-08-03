package database

import (
	"database/sql"

	"github.com/uptrace/bun"
)

// Connection interface defines the database connection contract
type Connection interface {
	GetDB() *bun.DB
	Close() error
	Health() error
	Stats() sql.DBStats
}
