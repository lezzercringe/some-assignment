package payments

import "github.com/shopspring/decimal"

type Payment struct {
	Transaction  string          `json:"transaction"`
	OrderID      string          `json:"-"`
	RequestID    string          `json:"request_id"`
	Currency     string          `json:"currency"`
	Provider     string          `json:"provider"`
	PaymentDT    int64           `json:"payment_dt"`
	Bank         string          `json:"bank"`
	GoodsTotal   int             `json:"goods_total"`
	Amount       decimal.Decimal `json:"amount"`
	DeliveryCost decimal.Decimal `json:"delivery_cost"`
	CustomFee    decimal.Decimal `json:"custom_fee"`
}
