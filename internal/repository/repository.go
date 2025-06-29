package repository

import (
	"context"
	"firstmod/internal/models"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	log  *slog.Logger
	conn *pgxpool.Pool
}

func New(log *slog.Logger, address string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Error("failed to ping database", "error", err)
		return nil, err
	}

	log.Info("successfully connected to database", "address", address)

	return &DB{
		log:  log,
		conn: pool,
	}, nil
}

func (db *DB) Add(ctx context.Context, order models.Order) error {
	return nil
}

func (db *DB) GetInfo(ctx context.Context, orderID int) (models.Order, error) {
	return models.Order{}, nil
}

func (db *DB) Delete(ctx context.Context, orderID int) error {
	return nil
}
