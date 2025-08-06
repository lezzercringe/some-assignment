package postgres

import (
	"context"
	"order-persistor/internal/orders"
	"order-persistor/internal/postgres/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ orders.Repository = &OrdersRepository{}

type OrdersRepository struct {
	ItemsDAO    *ItemsDAO
	PaymentsDAO *PaymentsDAO
	Pool        *pgxpool.Pool
}

func (r *OrdersRepository) Create(ctx context.Context, order *orders.Order) (*orders.Order, error) {
	inserted := &orders.Order{}

	err := withTx(ctx, r.Pool, func(ctx context.Context) error {
		var err error
		inserted, err = r.insertOrder(ctx, order)
		if err != nil {
			return err
		}

		inserted.Items, err = r.insertItems(ctx, order.ID, order.Items)
		if err != nil {
			return err
		}

		if order.Payment != nil {
			inserted.Payment, err = r.PaymentsDAO.Create(ctx, order.ID, order.Payment)
			if err != nil {
				return err
			}
		}

		return err
	})

	if err != nil {
		return nil, describeError(err)
	}

	return inserted, nil
}

func (r *OrdersRepository) GetByID(ctx context.Context, id string) (*orders.Order, error) {
	var order orders.Order

	err := withTx(ctx, r.Pool, func(ctx context.Context) error {
		tx := ctx.Value(txKey{}).(pgx.Tx)
		var err error
		dto, err := sqlc.New(tx).GetOrderByID(ctx, id)
		if err != nil {
			return err
		}

		order = *mapDtoToOrder(dto)

		order.Items, err = r.ItemsDAO.GetByOrderID(ctx, id)
		if err != nil {
			return err
		}

		payment, err := r.PaymentsDAO.GetByOrderID(ctx, id)
		if err != nil {
			return err
		}

		order.Payment = payment
		return nil
	})

	if err != nil {
		return nil, describeError(err)
	}

	return &order, nil
}

func (r *OrdersRepository) insertOrder(ctx context.Context, o *orders.Order) (*orders.Order, error) {
	executor := extractExecutor(ctx, r.Pool)
	inserted, err := sqlc.New(executor).CreateOrder(ctx, sqlc.CreateOrderParams{
		ID:                o.ID,
		TrackNumber:       o.TrackNumber,
		Entry:             o.Entry,
		Locale:            o.Locale,
		InternalSignature: o.Signature,
		CustomerID:        o.CustomerID,
		DeliveryService:   o.DeliveryService,
		Shardkey:          o.ShardKey,
		SmID:              int32(o.SMID),
		DateCreated:       o.CreatedAt,
		OofShard:          o.OOFShard,
		DeliveryName:      o.Delivery.Name,
		DeliveryPhone:     o.Delivery.Phone,
		DeliveryZip:       o.Delivery.Zip,
		DeliveryAddress:   o.Delivery.Address,
		DeliveryRegion:    o.Delivery.Region,
		DeliveryEmail:     o.Delivery.Email,
	})

	if err != nil {
		return nil, err
	}

	return mapDtoToOrder(inserted), nil
}

func (r *OrdersRepository) insertItems(ctx context.Context, orderID string, items []orders.Item) ([]orders.Item, error) {
	inserted := make([]orders.Item, 0, len(items))
	for _, item := range items {
		item, err := r.ItemsDAO.Create(ctx, orderID, &item)
		if err != nil {
			return nil, err
		}

		inserted = append(inserted, *item)
	}

	return inserted, nil
}

func mapDtoToOrder(o sqlc.Order) *orders.Order {
	return &orders.Order{
		ID:          o.ID,
		TrackNumber: o.TrackNumber,
		Entry:       o.Entry,
		Delivery: orders.Delivery{
			Name:    o.DeliveryName,
			Phone:   o.DeliveryPhone,
			Zip:     o.DeliveryZip,
			City:    o.DeliveryCity,
			Address: o.DeliveryAddress,
			Region:  o.DeliveryRegion,
			Email:   o.DeliveryEmail,
		},
		Locale:          o.Locale,
		Signature:       o.InternalSignature,
		CustomerID:      o.CustomerID,
		DeliveryService: o.DeliveryService,
		ShardKey:        o.Shardkey,
		SMID:            int(o.SmID),
		CreatedAt:       o.DateCreated,
		OOFShard:        o.OofShard,
	}
}
