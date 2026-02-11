package shared

import (
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Database struct {
	*sqlx.DB
	logger *zap.Logger
}

func NewDatabase(databaseURL string, logger *zap.Logger) (*Database, error) {
	// Normaliza URL: pgx5:// -> postgres://
	connStr := strings.Replace(databaseURL, "pgx5://", "postgres://", 1)
	db, err := sqlx.Connect("pgx", connStr)

	if err != nil {
		logger.Error("failed to connect to database", zap.Error(err))
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	return &Database{
		DB:     db,
		logger: logger,
	}, nil
}

func (d *Database) Close() error {
	d.logger.Info("closing database connection")
	return d.DB.Close()
}

func (d *Database) HealthCheck() error {
	return d.DB.Ping()
}
