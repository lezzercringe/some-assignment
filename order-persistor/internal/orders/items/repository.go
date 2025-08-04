package items

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Item, error)
	Create(ctx context.Context, i *Item) error
	Update(ctx context.Context, i *Item) error
	Delete(ctx context.Context, id string)
}
