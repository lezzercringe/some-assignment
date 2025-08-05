package payments

import "context"

type Repository interface {
	Create(ctx context.Context, p *Payment) (*Payment, error)
	GetByOrderID(ctx context.Context, orderID string) (*Payment, error)
}
