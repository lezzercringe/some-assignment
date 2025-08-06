package orders

import "github.com/shopspring/decimal"

type Payment struct {
	Transaction  string          `json:"transaction" validate:"required"`
	RequestID    string          `json:"request_id"`
	Currency     string          `json:"currency" validate:"required"`
	Provider     string          `json:"provider" validate:"required"`
	PaymentDT    int64           `json:"payment_dt" validate:"required"`
	Bank         string          `json:"bank" validate:"required"`
	GoodsTotal   int             `json:"goods_total" validate:"required"`
	Amount       decimal.Decimal `json:"amount" validate:"required"`
	DeliveryCost decimal.Decimal `json:"delivery_cost" validate:"required"`
	CustomFee    decimal.Decimal `json:"custom_fee" validate:"required"`
}
