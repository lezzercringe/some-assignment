package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type Executor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
}

func extractExecutor(ctx context.Context, pool *pgxpool.Pool) Executor {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}

	return pool
}

func withTx(ctx context.Context, pool *pgxpool.Pool, fn func(ctx context.Context) error) error {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	if !ok {
		var err error
		tx, err = pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("error starting new transaction: %w", err)
		}

		ctx = context.WithValue(ctx, txKey{}, tx)
	}

	if err := fn(ctx); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return fmt.Errorf("error rolling back tx: %w", err)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commiting transaction: %w", err)
	}

	return nil
}
