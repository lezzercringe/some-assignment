package postgres

import (
	"context"
	"order-persistor/internal/orders"
	"order-persistor/internal/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentsDAO struct {
	Pool *pgxpool.Pool
}

func (r *PaymentsDAO) GetByOrderID(ctx context.Context, orderID string) (*orders.Payment, error) {
	exec := extractExecutor(ctx, r.Pool)
	dto, err := sqlc.New(exec).GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return mapDtoToPayment(dto), nil
}

func (r *PaymentsDAO) Create(ctx context.Context, orderID string, p *orders.Payment) (*orders.Payment, error) {
	exec := extractExecutor(ctx, r.Pool)
	dto, err := sqlc.New(exec).CreatePayment(ctx, sqlc.CreatePaymentParams{
		Transaction:  p.Transaction,
		OrderID:      orderID,
		RequestID:    p.RequestID,
		Currency:     p.Currency,
		Provider:     p.Provider,
		Amount:       p.Amount,
		PaymentDt:    p.PaymentDT,
		Bank:         p.Bank,
		DeliveryCost: p.DeliveryCost,
		GoodsTotal:   int32(p.GoodsTotal),
		CustomFee:    p.CustomFee,
	})

	if err != nil {
		return nil, err
	}

	return mapDtoToPayment(dto), nil
}

func mapDtoToPayment(p sqlc.Payment) *orders.Payment {
	return &orders.Payment{
		Transaction:  p.Transaction,
		RequestID:    p.RequestID,
		Currency:     p.Currency,
		Provider:     p.Provider,
		PaymentDT:    p.PaymentDt,
		Bank:         p.Bank,
		GoodsTotal:   int(p.GoodsTotal),
		Amount:       p.Amount,
		DeliveryCost: p.DeliveryCost,
		CustomFee:    p.CustomFee,
	}
}
