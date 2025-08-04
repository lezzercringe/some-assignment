package postgres

import (
	"context"
	"order-persistor/internal/orders/deliveries"
)

var _ deliveries.Repository = &DeliveriesRepository{}

type DeliveriesRepository struct{}

func (*DeliveriesRepository) Create(ctx context.Context, d *deliveries.Delivery) error {
	panic("unimplemented")
}

func (d *DeliveriesRepository) Delete(ctx context.Context, id string) {
	panic("unimplemented")
}

func (d *DeliveriesRepository) GetByID(ctx context.Context, id string) (*deliveries.Delivery, error) {
	panic("unimplemented")
}

func (*DeliveriesRepository) Update(ctx context.Context, d *deliveries.Delivery) error {
	panic("unimplemented")
}
