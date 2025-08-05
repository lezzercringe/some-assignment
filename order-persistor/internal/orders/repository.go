package orders

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, o *Order) (*Order, error)
}
