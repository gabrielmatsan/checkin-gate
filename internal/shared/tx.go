package shared

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// DBTX é uma interface que tanto *sqlx.DB quanto *sqlx.Tx implementam
type DBTX interface {
	sqlx.ExtContext
	sqlx.PreparerContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// TxManager gerencia transações de banco de dados
type TxManager struct {
	db *sqlx.DB
}

func NewTxManager(db *sqlx.DB) *TxManager {
	return &TxManager{db: db}
}

// WithTx executa uma função dentro de uma transação
// Se a função retornar erro, faz rollback; caso contrário, faz commit
func (tm *TxManager) WithTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := tm.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if err := fn(tx); err != nil {
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

// WithTxResult executa uma função dentro de uma transação e retorna um resultado
func WithTxResult[T any](tm *TxManager, ctx context.Context, fn func(tx *sqlx.Tx) (T, error)) (T, error) {
	var result T

	tx, err := tm.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to begin transaction: %w", err)
	}

	result, err = fn(tx)
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
