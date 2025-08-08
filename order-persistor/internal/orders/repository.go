package orders

import (
	"context"
	"errors"
)

var (
	ErrInternalFailure = errors.New("internal db failure")
	ErrNotFound        = errors.New("not found")
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*Order, error)
	ListRecent(ctx context.Context, n int) ([]Order, error)
	Create(ctx context.Context, o *Order) (*Order, error)
}
