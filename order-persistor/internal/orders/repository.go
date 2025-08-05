package orders

import (
	"context"
	"errors"
)

var ErrInternalFailure = errors.New("internal db failure")

type Repository interface {
	GetByID(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, o *Order) (*Order, error)
}
