package orders

import (
	"order-persistor/internal/orders/deliveries"
	"order-persistor/internal/orders/items"
	"order-persistor/internal/orders/payments"
	"time"
)

type Order struct {
	ID              string              `json:"order_uid"`
	TrackNumber     string              `json:"track_number"`
	Entry           string              `json:"entry"`
	Delivery        deliveries.Delivery `json:"delivery"`
	Payment         payments.Payment    `json:"payment"`
	Items           []items.Item        `json:"items"`
	Locale          string              `json:"locale"`
	Signature       string              `json:"internal_signature"`
	CustomerID      string              `json:"customer_id"`
	DeliveryService string              `json:"delivery_service"`
	ShardKey        string              `json:"shardkey"`
	SMID            int                 `json:"sm_id"`
	CreatedAt       time.Time           `json:"date_created"`
	OOFShard        string              `json:"oof_shard"`
}
