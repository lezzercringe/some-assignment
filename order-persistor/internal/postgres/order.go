package postgres

import (
	"context"
	"order-persistor/internal/orders"
	"order-persistor/internal/orders/deliveries"
	"order-persistor/internal/orders/items"
)

var _ orders.Repository = &OrdersRepository{}

type OrdersRepository struct {
	ItemRepository items.Repository
	Deliveries     deliveries.Repository
}

func (*OrdersRepository) Create(ctx context.Context, o *orders.Order) error {
	panic("unimplemented")
}

func (o *OrdersRepository) Delete(ctx context.Context, id string) {
	panic("unimplemented")
}

func (o *OrdersRepository) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	panic("unimplemented")
}

func (*OrdersRepository) Update(ctx context.Context, o *orders.Order) error {
	panic("unimplemented")
}
