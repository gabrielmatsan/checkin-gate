package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gabrielmatsan/checkin-gate/internal/events/domain/repository"
	"github.com/jmoiron/sqlx"
)

type PostgresTransactionProvider struct {
	db *sqlx.DB
}

func NewPostgresTransactionProvider(db *sqlx.DB) *PostgresTransactionProvider {
	return &PostgresTransactionProvider{db: db}
}

func (p *PostgresTransactionProvider) Transact(ctx context.Context, fn func(repos repository.Repositories) error) error {
	tx, err := p.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	repos := repository.Repositories{
		Events:     NewPostgresEventRepository(tx),
		Activities: NewPostgresActivityRepository(tx),
		CheckIns:   NewPostgresCheckInRepository(tx),
	}

	if err := fn(repos); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TransactResult executa uma função dentro de uma transação e retorna um resultado
func TransactResult[T any](p *PostgresTransactionProvider, ctx context.Context, fn func(repos repository.Repositories) (T, error)) (T, error) {
	var result T

	tx, err := p.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}

	repos := repository.Repositories{
		Events:     NewPostgresEventRepository(tx),
		Activities: NewPostgresActivityRepository(tx),
		CheckIns:   NewPostgresCheckInRepository(tx),
	}

	result, err = fn(repos)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return result, fmt.Errorf("tx error: %v, rollback error: %w", err, rbErr)
		}
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}
