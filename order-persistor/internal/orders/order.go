package orders

import "time"

type Order struct {
	ID              string    `json:"order_uid" validate:"required"`
	TrackNumber     string    `json:"track_number" validate:"required"`
	Entry           string    `json:"entry" validate:"required"`
	Delivery        Delivery  `json:"delivery"`
	Payment         *Payment  `json:"payment"`
	Items           []Item    `json:"items"`
	Locale          string    `json:"locale" validate:"required"`
	Signature       string    `json:"internal_signature"`
	CustomerID      string    `json:"customer_id" validate:"required"`
	DeliveryService string    `json:"delivery_service" validate:"required"`
	ShardKey        string    `json:"shardkey" validate:"required"`
	SMID            int       `json:"sm_id" validate:"required"`
	CreatedAt       time.Time `json:"date_created" validate:"required"`
	OOFShard        string    `json:"oof_shard" validate:"required,numeric"`
}
