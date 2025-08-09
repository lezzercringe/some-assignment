package producer

import (
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
)

func GenerateOrder() Order {
	return Order{
		ID:          faker.UUIDHyphenated(),
		TrackNumber: faker.Word(),
		Entry:       faker.Word(),
		Delivery: Delivery{
			Name:    faker.Name(),
			Phone:   faker.Phonenumber(),
			Zip:     faker.Word(),
			City:    faker.Word(),
			Address: faker.Word(),
			Region:  faker.Word(),
			Email:   faker.Email(),
		},
		Payment: Payment{
			Transaction:  faker.UUIDHyphenated(),
			RequestID:    faker.UUIDHyphenated(),
			Currency:     "USD",
			Provider:     faker.Word(),
			Amount:       rand.Intn(10000),
			PaymentDT:    time.Now().Unix(),
			Bank:         faker.Word(),
			DeliveryCost: rand.Intn(2000),
			GoodsTotal:   rand.Intn(3000),
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:      rand.Intn(9999999),
				TrackNumber: faker.Word(),
				Price:       453,
				RID:         faker.UUIDHyphenated(),
				Name:        faker.Word(),
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        rand.Intn(9999999),
				Brand:       faker.Word(),
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        faker.Username(),
		DeliveryService:   faker.Word(),
		ShardKey:          "9",
		SmID:              rand.Intn(999),
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
}
