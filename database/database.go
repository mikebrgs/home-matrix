package database

import (
	"context"

	"github.com/mikebrgs/home-matrix/models/pot"
)

// Database defines the interface for database operations
type Database interface {
	Connect(ctx context.Context, dsn string) error
	Close(ctx context.Context) error
	EnsurePotHealthTableExists(ctx context.Context) error
	EnsurePotStatusTableExists(ctx context.Context) error
	InsertPotHealthData(ctx context.Context, data pot.PotHealth) error
	InsertPotStatusData(ctx context.Context, data pot.PotStatus) error
	GetPotHealthData(ctx context.Context) ([]pot.PotHealth, error)
	GetPotStatusData(ctx context.Context) ([]pot.PotStatus, error)
}
