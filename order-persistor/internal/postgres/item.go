package postgres

import (
	"context"
	"order-persistor/internal/orders/items"
)

var _ items.Repository = &ItemRepository{}

type ItemRepository struct{}

func (*ItemRepository) Create(ctx context.Context, i *items.Item) error {
	panic("unimplemented")
}

func (i *ItemRepository) Delete(ctx context.Context, id string) {
	panic("unimplemented")
}

func (i *ItemRepository) GetByID(ctx context.Context, id string) (*items.Item, error) {
	panic("unimplemented")
}

func (*ItemRepository) Update(ctx context.Context, i *items.Item) error {
	panic("unimplemented")
}
