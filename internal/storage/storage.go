package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type storage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

var schema = `
CREATE TABLE IF NOT EXISTS "accounts" (
	account_id text UNIQUE PRIMARY KEY,
	funds numeric default 0,
	on_hold numeric default 0
	);

CREATE TABLE IF NOT EXISTS "orders" (
	order_id bigint UNIQUE PRIMARY KEY,
	service_id text,
	account_id text,
	amount numeric,
	creation_date timestamptz,
	processed_date timestamptz,
	processed bool default false
	);

CREATE TABLE IF NOT EXISTS "reports" (
    report_id text UNIQUE PRIMARY KEY,
    report text NOT NULL
);

ALTER TABLE orders ADD FOREIGN KEY (account_id) REFERENCES accounts (account_id);
`

// New creates a new instance of database layer and migrates it
func New(path string, logger *zap.Logger) *storage {
	// Wait until database initialize in container
	time.Sleep(time.Second * 2)

	config, err := pgxpool.ParseConfig(path)
	if err != nil {
		logger.Error("Unable to parse config", zap.Error(err))
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("Unable to connect to database", zap.Error(err))
	}

	return &storage{
		pool:   conn,
		logger: logger,
	}
}

// CreateSchema executes needed schema
func (s storage) CreateSchema() {
	_, err := s.pool.Exec(context.Background(), schema)
	if err != nil {
		s.logger.Error("Error occurred while executing schema", zap.Error(err))
	}

	s.logger.Info("Launched with pgx. Database created.")
}
