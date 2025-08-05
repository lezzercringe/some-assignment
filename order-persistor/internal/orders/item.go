package orders

import "github.com/shopspring/decimal"

type Item struct {
	CHRTID      int             `json:"chrt_id"`
	TrackNumber string          `json:"track_number"`
	RID         string          `json:"rid"`
	Name        string          `json:"name"`
	Size        string          `json:"size"`
	NMID        int             `json:"nm_id"`
	Brand       string          `json:"brand"`
	Status      int             `json:"status"`
	Price       decimal.Decimal `json:"price"`
	Sale        decimal.Decimal `json:"sale"`
	TotalPrice  decimal.Decimal `json:"total_price"`
}
