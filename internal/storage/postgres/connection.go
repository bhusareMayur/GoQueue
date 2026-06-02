package postgres

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	// Tweaked for Horizontal Scaling
	config.MaxConns = 80 
	config.MinConns = 5 // Reduced from 20 to allow more workers to share the DB limit safely

	slog.Info("configuring postgres pool", "max_conns", config.MaxConns, "min_conns", config.MinConns)

	return pgxpool.NewWithConfig(
		context.Background(),
		config,
	)
}