package payments

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Payment, error)
	Create(ctx context.Context, p *Payment) error
	Update(ctx context.Context, p *Payment) error
	Delete(ctx context.Context, id string)
}
