package inmemory

import (
	"context"
	"fmt"
	"log/slog"
	"order-persistor/internal/config"
	"order-persistor/internal/orders"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

var _ orders.Repository = &OrdersCache{}

type OrdersCache struct {
	decoratee orders.Repository
	size      int
	lru       *lru.Cache[string, *orders.Order]
	logger    *slog.Logger
}

// NewOrdersCache creates a ready-to-use LRU cache.
// It does not bootstrap or pre-fill the cache, even if pre-filling is enabled in the configuration,
// because it is a potentially long-running operation.
// To manually trigger pre-filling after creation, use the Prefill method.
func NewOrdersCache(cfg config.Cache, decoratee orders.Repository, logger *slog.Logger) (*OrdersCache, error) {
	lru, err := lru.New[string, *orders.Order](cfg.Size)
	if err != nil {
		return nil, fmt.Errorf("could not create lru: %w", err)
	}

	return &OrdersCache{
		decoratee: decoratee,
		lru:       lru,
		logger:    logger,
		size:      cfg.Size,
	}, nil
}

// Prefill fills up the cache with the most fresh orders.
func (c *OrdersCache) Prefill(ctx context.Context) error {
	now := time.Now()
	qtyLoaded, err := c.load(ctx, c.size)
	took := time.Since(now)
	if err != nil {
		return err
	}

	c.logger.Info("cache: prefilled", "orders_loaded", qtyLoaded, "took", took.String())
	return nil
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

func (c *OrdersCache) ListRecent(ctx context.Context, n int) ([]orders.Order, error) {
	return c.decoratee.ListRecent(ctx, n)
}

func (c *OrdersCache) load(ctx context.Context, n int) (loaded int, err error) {
	start := time.Now()

	recents, err := c.decoratee.ListRecent(ctx, n)
	if err != nil {
		return 0, fmt.Errorf("could not load recent orders: %w", err)
	}

	for _, order := range recents {
		c.lru.Add(order.ID, &order)
	}

	duration := time.Since(start)
	c.logger.Info("cache: loaded up", "time", duration.String(), "entries", len(recents))
	return len(recents), nil
}
