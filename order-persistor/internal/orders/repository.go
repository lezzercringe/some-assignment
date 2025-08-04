package orders

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Order, error)
	Create(ctx context.Context, o *Order) error
	Update(ctx context.Context, o *Order) error
	Delete(ctx context.Context, id string)
}
