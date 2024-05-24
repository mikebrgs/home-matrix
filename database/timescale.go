package database

import (
	"context"

	"github.com/mikebrgs/home-matrix/models/pot"

	"github.com/jackc/pgx/v5"
)

type PostgresDB struct {
	conn *pgx.Conn
}

func (db *PostgresDB) Connect(ctx context.Context, dsn string) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}
	db.conn = conn
	return nil
}

func (db *PostgresDB) Close(ctx context.Context) error {
	if db.conn != nil {
		return db.conn.Close(ctx)
	}
	return nil
}

// func (db *PostgresDB) ensurePotHealthTableExists(ctx context.Context) {
// 	createTableQuery := `
// 	CREATE TABLE IF NOT EXISTS pot_health (
// 		ts BIGINT,
// 		temperature DOUBLE PRECISION,
// 		humidity DOUBLE PRECISION,
// 		pressure DOUBLE PRECISION,
// 		moisture DOUBLE PRECISION,
// 		light DOUBLE PRECISION
// 	);
// 	`
// 	_, err := db.conn.Exec(context.Background(), createTableQuery)
// 	if err != nil {
// 		slog.Error(err.Error())
// 	} else {
// 		slog.Info("Ensured table exists or created it if not.")
// 	}
// }

// // If the table exists, skip otherwise create a PotStatus table. Takes in
// // context, returns nothing.
// func (db *PostgresDB) ensurePotStatusTableExists(ctx context.Context) {
// 	createTableQuery := `
// 	CREATE TABLE IF NOT EXISTS pot_status (
// 		time TIMESTAMPTZ,
// 		pot_id TEXT,
// 		temperature DOUBLE PRECISION,
// 		humidity DOUBLE PRECISION,
// 		pressure DOUBLE PRECISION,
// 		moisture DOUBLE PRECISION,
// 		light DOUBLE PRECISION
// 	);
// 	`
// 	_, err := db.conn.Exec(ctx, createTableQuery)
// 	if err != nil {
// 		slog.Error(err.Error())
// 	} else {
// 		slog.Info("Ensured table exists or created it doesn't.")
// 	}
// }

// func (db *PostgresDB) ensurePotStatusTableIsHypertable(ctx context.Context) {
// 	createHyperTableQuery := `
// 		SELECT create_hypertable('pot_status', by_range('time'), if_not_exists => TRUE);
// 	`
// 	_, err := db.conn.Exec(ctx, createHyperTableQuery)
// 	if err != nil {
// 		slog.Error(err.Error())
// 	} else {
// 		slog.Info("Ensured table exists or created it doesn't.")
// 	}
// }

func (db *PostgresDB) InsertPotHealthData(ctx context.Context, data pot.PotHealth) error {
	_, err := db.conn.Exec(ctx,
		"INSERT INTO pot_health (ts, temperature, humidity, pressure, moisture, light) VALUES ($1, $2, $3, $4, $5, $6)",
		data.TS, data.Temperature, data.Humidity, data.Pressure, data.Moisture, data.Light)
	return err
}

func (db *PostgresDB) InsertPotStatusData(ctx context.Context, data pot.PotStatus) error {
	_, err := db.conn.Exec(ctx,
		"INSERT INTO pot_status (ts, battery, cpu, memory, storage, temperature) VALUES ($1, $2, $3, $4, $5, $6)",
		data.TS, data.Battery, data.CPU, data.Memory, data.Storage, data.Temperature)
	return err
}
