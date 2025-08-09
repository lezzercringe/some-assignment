package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"order-persistor/internal/config"
	"order-persistor/internal/mocks"
	"order-persistor/internal/orders"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"
)

func newTestConsumer(client Client, repository orders.Repository) *OrdersConsumer {
	return &OrdersConsumer{
		client: client,
		cfg: config.KafkaConsumer{
			Topic:              "test-topic",
			ReadTimeout:        1 * time.Second,
			ProcessTimeout:     1 * time.Second,
			ReadFailureBackoff: 1 * time.Second,
		},
		ordersRepository: repository,
		logger:           slog.New(slog.DiscardHandler),
		shutdownOnce:     sync.Once{},
	}
}

func TestOrdersConsumer_handleMessage(t *testing.T) {
	t.Parallel()

	t.Run("internal error is propagated", func(t *testing.T) {
		wantError := orders.ErrInternalFailure

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := mocks.NewMockClient(ctrl)
		consumer := newTestConsumer(client, repository)

		repository.EXPECT().
			Create(gomock.Any(), gomock.Eq(&validOrder)).
			Return(nil, wantError).
			Times(1)

		jsonEncoded, _ := json.Marshal(validOrder)

		err := consumer.handleMessage(
			context.Background(),
			jsonEncoded,
		)

		if err == nil {
			t.Fatalf("consumer did not propagate internal error")
		}

		if errors.Is(err, errMalformedOrder) {
			t.Fatalf("error wraps errMalformedOrder")
		}

		if !errors.Is(err, wantError) {
			t.Fatalf("error returned from consumer does not wrap expected error")
		}
	})

	t.Run("bad json - returns malformed message error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := mocks.NewMockClient(ctrl)
		consumer := newTestConsumer(client, repository)

		err := consumer.handleMessage(
			context.Background(),
			[]byte("some bad json"),
		)

		if err == nil || !errors.Is(err, errMalformedOrder) {
			t.Fatalf("did not return malformed message error")
		}
	})

	t.Run("validation failure - returns malformed message error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := mocks.NewMockClient(ctrl)
		consumer := newTestConsumer(client, repository)

		jsonEncoded, _ := json.Marshal(invalidOrder)

		err := consumer.handleMessage(
			context.Background(),
			jsonEncoded,
		)

		if err == nil || !errors.Is(err, errMalformedOrder) {
			t.Fatalf("did not return malformed message error")
		}
	})

	t.Run("context cancellation - does return ctx error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repository := mocks.NewMockRepository(ctrl)
		client := mocks.NewMockClient(ctrl)
		consumer := newTestConsumer(client, repository)

		ctx, cancel := context.WithCancel(context.Background())

		repository.EXPECT().
			Create(gomock.Any(), gomock.Eq(&validOrder)).
			DoAndReturn(func(createCtx context.Context, o *orders.Order) (*orders.Order, error) {
				cancel()
				<-createCtx.Done()
				return nil, createCtx.Err()
			}).Times(1)

		jsonEncoded, _ := json.Marshal(validOrder)

		err := consumer.handleMessage(
			ctx,
			jsonEncoded,
		)

		if err == nil {
			t.Fatalf("consumer did not return context error")
		}

		if !errors.Is(err, ctx.Err()) {
			t.Fatalf("the returned error does not wrap context")
		}
	})
}

var validOrder = orders.Order{
	ID:          "order-0001",
	TrackNumber: "TRK123456789",
	Entry:       "ENTRY-ABC-1",

	Delivery: orders.Delivery{
		Name:    "Jane Doe",
		Phone:   "+358401234567",
		Zip:     "00100",
		City:    "Helsinki",
		Address: "Testintie 1 A 2",
		Region:  "Uusimaa",
		Email:   "jane.doe@example.com",
	},

	Payment: &orders.Payment{
		Transaction: "txn_20250809_01",
		RequestID:   "req-2025-08-09-01",
		Currency:    "EUR",
		Provider:    "stripe",
		PaymentDT:   time.Now().Unix(),

		Bank:         "Test Bank Oy",
		GoodsTotal:   249,
		DeliveryCost: decimal.NewFromFloat(9.99),
		CustomFee:    decimal.NewFromFloat(0.00),
		Amount:       decimal.NewFromFloat(249.49).Add(decimal.NewFromFloat(9.99)),
	},

	Items: []orders.Item{
		{
			CHRTID:      1001,
			TrackNumber: "ITM-TRK-1",
			RID:         "rid-1",
			Name:        "Comfort Sneakers",
			Size:        "42",
			NMID:        5001,
			Brand:       "SneakerCo",
			Status:      1,
			Price:       decimal.NewFromFloat(199.99),
			Sale:        decimal.NewFromFloat(0.00),
			TotalPrice:  decimal.NewFromFloat(199.99),
		},
		{
			CHRTID:      1002,
			TrackNumber: "ITM-TRK-2",
			RID:         "rid-2",
			Name:        "Everyday Socks (3-pack)",
			Size:        "L",
			NMID:        5002,
			Brand:       "SockMakers",
			Status:      1,
			Price:       decimal.NewFromFloat(49.50),
			Sale:        decimal.NewFromFloat(0.00),
			TotalPrice:  decimal.NewFromFloat(49.50),
		},
	},

	Locale:          "en-US",
	Signature:       "sig-example-base64==",
	CustomerID:      "cust-007",
	DeliveryService: "DHL",
	ShardKey:        "shard-1",
	SMID:            42,
	CreatedAt:       time.Date(2021, 11, 14, 8, 27, 53, 123456000, time.UTC),
	OOFShard:        "1",
}

// no id stated
var invalidOrder = orders.Order{
	TrackNumber: "TRK123456789",
	Entry:       "ENTRY-ABC-1",

	Delivery: orders.Delivery{
		Name:    "Jane Doe",
		Phone:   "+358401234567",
		Zip:     "00100",
		City:    "Helsinki",
		Address: "Testintie 1 A 2",
		Region:  "Uusimaa",
		Email:   "jane.doe@example.com",
	},

	Payment: &orders.Payment{
		Transaction: "txn_20250809_01",
		RequestID:   "req-2025-08-09-01",
		Currency:    "EUR",
		Provider:    "stripe",
		PaymentDT:   time.Now().Unix(),

		Bank:         "Test Bank Oy",
		GoodsTotal:   249,
		DeliveryCost: decimal.NewFromFloat(9.99),
		CustomFee:    decimal.NewFromFloat(0.00),
		Amount:       decimal.NewFromFloat(249.49).Add(decimal.NewFromFloat(9.99)),
	},

	Items: []orders.Item{
		{
			CHRTID:      1001,
			TrackNumber: "ITM-TRK-1",
			RID:         "rid-1",
			Name:        "Comfort Sneakers",
			Size:        "42",
			NMID:        5001,
			Brand:       "SneakerCo",
			Status:      1,
			Price:       decimal.NewFromFloat(199.99),
			Sale:        decimal.NewFromFloat(0.00),
			TotalPrice:  decimal.NewFromFloat(199.99),
		},
		{
			CHRTID:      1002,
			TrackNumber: "ITM-TRK-2",
			RID:         "rid-2",
			Name:        "Everyday Socks (3-pack)",
			Size:        "L",
			NMID:        5002,
			Brand:       "SockMakers",
			Status:      1,
			Price:       decimal.NewFromFloat(49.50),
			Sale:        decimal.NewFromFloat(0.00),
			TotalPrice:  decimal.NewFromFloat(49.50),
		},
	},

	Locale:          "en-US",
	Signature:       "sig-example-base64==",
	CustomerID:      "cust-007",
	DeliveryService: "DHL",
	ShardKey:        "shard-1",
	SMID:            42,
	CreatedAt:       time.Date(2021, 11, 14, 8, 27, 53, 123456000, time.UTC),
	OOFShard:        "1",
}
