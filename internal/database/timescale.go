package database

import (
	"context"

	"github.com/mikebrgs/home-matrix/internal/models/pot"

	"github.com/jackc/pgx/v5"
)

func NewTimescaleDBClient(connStr string) (*TimescaleDB, error) {
	ctx := context.Background()
	db := TimescaleDB{
		&pgx.Conn{},
		&ctx,
	}
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	db.conn = conn
	return &db, nil
}

type TimescaleDB struct {
	conn *pgx.Conn
	ctx  *context.Context
}

func (db *TimescaleDB) Connect() error {
	return nil
}

func (db *TimescaleDB) Close() error {
	if db.conn != nil {
		return db.conn.Close(*db.ctx)
	}
	return nil
}

func (db *TimescaleDB) InsertPotHealthData(data pot.PotHealth) error {
	_, err := db.conn.Exec(*db.ctx,
		"INSERT INTO pot_health (time, device_id, temperature, humidity, pressure, moisture, light) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		data.TS, data.DeviceId, data.Temperature, data.Humidity, data.Pressure, data.Moisture, data.Light)
	return err
}

func (db *TimescaleDB) InsertPotStatusData(data pot.PotStatus) error {
	_, err := db.conn.Exec(*db.ctx,
		"INSERT INTO pot_status (time, device_id, battery, cpu, memory, storage, temperature) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		data.TS, data.DeviceId, data.Battery, data.CPU, data.Memory, data.Storage, data.Temperature)
	return err
}
