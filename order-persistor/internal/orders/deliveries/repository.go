package deliveries

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Delivery, error)
	Create(ctx context.Context, d *Delivery) error
	Update(ctx context.Context, d *Delivery) error
	Delete(ctx context.Context, id string)
}
