package database

import (
	"context"

	"github.com/mikebrgs/home-matrix/internal/models/pot"
)

// Database defines the interface for database operations
type Database interface {
	Connect(dsn string) error
	Close(ctx context.Context) error
	InsertPotHealthData(data pot.PotHealth) error
	InsertPotStatusData(data pot.PotStatus) error
}
