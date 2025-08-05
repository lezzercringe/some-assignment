package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"order-persistor/internal/config"
	"order-persistor/internal/orders"

	lru "github.com/hashicorp/golang-lru/v2"
)

var _ orders.Repository = &OrdersCache{}

type OrdersCache struct {
	decoratee orders.Repository
	lru       *lru.Cache[string, *orders.Order]
	logger    *slog.Logger
}

func NewOrdersCache(cfg config.Cache, decoratee orders.Repository, logger *slog.Logger) (*OrdersCache, error) {
	lru, err := lru.New[string, *orders.Order](cfg.Size)
	if err != nil {
		return nil, fmt.Errorf("could not create lru: %w", err)
	}

	return &OrdersCache{
		decoratee: decoratee,
		lru:       lru,
		logger:    logger,
	}, nil
}

func (c *OrdersCache) Create(ctx context.Context, o *orders.Order) (*orders.Order, error) {
	inserted, err := c.decoratee.Create(ctx, o)
	if err != nil {
		return nil, err
	}

	c.lru.Add(inserted.ID, inserted)
	return inserted, nil
}

func (c *OrdersCache) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	order, hit := c.lru.Get(id)
	if hit {
		return order, nil
	}

	c.logger.DebugContext(ctx,
		"order cache miss",
		slog.String("order_id", id),
		slog.String("op", "GetByID"),
	)

	order, err := c.decoratee.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.lru.Add(id, order)
	return order, nil
}
