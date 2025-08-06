package postgres

import (
	"context"
	"order-persistor/internal/orders"
	"order-persistor/internal/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemsDAO struct {
	Pool *pgxpool.Pool
}

func (r *ItemsDAO) Create(ctx context.Context, orderID string, i *orders.Item) (*orders.Item, error) {
	exec := extractExecutor(ctx, r.Pool)
	item, err := sqlc.New(exec).CreateItem(ctx, sqlc.CreateItemParams{
		OrderID:     orderID,
		ChrtID:      int32(i.CHRTID),
		TrackNumber: i.TrackNumber,
		Price:       i.Price,
		Rid:         i.RID,
		Name:        i.Name,
		Sale:        i.Sale,
		Size:        i.Size,
		TotalPrice:  i.TotalPrice,
		NmID:        int32(i.NMID),
		Brand:       i.Brand,
		Status:      int32(i.Status),
	})

	if err != nil {
		return nil, err
	}

	return mapDtoToItem(item), nil
}

func (r *ItemsDAO) GetByOrderID(ctx context.Context, orderID string) ([]orders.Item, error) {
	exec := extractExecutor(ctx, r.Pool)
	dtos, err := sqlc.New(exec).GetItemsByOrderID(ctx, orderID)

	if err != nil {
		return nil, err
	}

	items := make([]orders.Item, 0, len(dtos))
	for _, dto := range dtos {
		items = append(items, *mapDtoToItem(dto))
	}

	return items, nil
}

func mapDtoToItem(i sqlc.Item) *orders.Item {
	return &orders.Item{
		CHRTID:      int(i.ChrtID),
		TrackNumber: i.TrackNumber,
		RID:         i.Rid,
		Name:        i.Name,
		Size:        i.Size,
		NMID:        int(i.NmID),
		Brand:       i.Brand,
		Status:      int(i.Status),
		Price:       i.Price,
		Sale:        i.Sale,
		TotalPrice:  i.TotalPrice,
	}

}
