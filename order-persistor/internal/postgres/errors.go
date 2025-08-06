package postgres

import (
	"errors"
	"fmt"
	"order-persistor/internal/orders"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func describeError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return orders.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("already exists (unique violation): %w", err)
		case "23514":
			return fmt.Errorf("validation failed (check violation): %w", err)
		case "23502":
			return fmt.Errorf("missing required field (not null violation): %w", err)
		}
	}

	return errors.Join(orders.ErrInternalFailure, err)
}
