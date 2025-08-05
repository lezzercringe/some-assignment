package postgres

import (
	"context"
	"order-persistor/internal/orders/payments"
	"order-persistor/internal/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ payments.Repository = &PaymentsRepository{}

type PaymentsRepository struct {
	Pool *pgxpool.Pool
}

func (r *PaymentsRepository) GetByOrderID(ctx context.Context, orderID string) (*payments.Payment, error) {
	exec := extractExecutor(ctx, r.Pool)
	dto, err := sqlc.New(exec).GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return mapDtoToPayment(dto), nil
}

func (r *PaymentsRepository) Create(ctx context.Context, p *payments.Payment) (*payments.Payment, error) {
	exec := extractExecutor(ctx, r.Pool)
	dto, err := sqlc.New(exec).CreatePayment(ctx, sqlc.CreatePaymentParams{
		Transaction:  p.Transaction,
		OrderID:      p.OrderID,
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

func mapDtoToPayment(p sqlc.Payment) *payments.Payment {
	return &payments.Payment{
		Transaction:  p.Transaction,
		OrderID:      p.OrderID,
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
