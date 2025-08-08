package inmemory

import (
	"context"
	"errors"
	"log/slog"
	"order-persistor/internal/config"
	"order-persistor/internal/mocks"
	"order-persistor/internal/orders"
	"reflect"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
)

func TestOrdersCache_GetByID(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.DiscardHandler)
	testOrder := &orders.Order{
		ID:        "someid",
		CreatedAt: time.Now(),
	}

	t.Run("cache miss", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rep := mocks.NewMockRepository(ctrl)

		rep.EXPECT().
			GetByID(gomock.Any(), gomock.Eq(testOrder.ID)).
			Return(testOrder, nil).
			Times(1)

		cache, err := NewOrdersCache(config.Cache{Size: 1}, rep, log)
		if err != nil {
			t.Fatalf("error creating cache: %v", err)
		}

		order, err := cache.GetByID(context.Background(), testOrder.ID)
		if err != nil {
			t.Fatalf("error while extracting existing order from cache: %v", err)
		}

		if !reflect.DeepEqual(testOrder, order) {
			t.Fatal("unexpected order extracted from cache")
		}
	})

	t.Run("cache hit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rep := mocks.NewMockRepository(ctrl)

		cache, err := NewOrdersCache(config.Cache{Size: 1}, rep, log)
		if err != nil {
			t.Fatalf("error creating cache: %v", err)
		}

		// pre-inserting value to ensure cache-hit
		cache.lru.Add(testOrder.ID, testOrder)

		order, err := cache.GetByID(context.Background(), testOrder.ID)
		if err != nil {
			t.Fatalf("error while extracting existing order from cache: %v", err)
		}

		if !reflect.DeepEqual(testOrder, order) {
			t.Fatal("unexpected order extracted from cache")
		}
	})
}

func TestOrdersCache_Create(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.DiscardHandler)
	testOrder := &orders.Order{
		ID:        "someid",
		CreatedAt: time.Now(),
	}

	t.Run("error gets propagated", func(t *testing.T) {
		decorateeErr := errors.New("some obscure error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rep := mocks.NewMockRepository(ctrl)

		rep.EXPECT().
			Create(gomock.Any(), gomock.Eq(testOrder)).
			Return(nil, decorateeErr).
			Times(1)

		cache, err := NewOrdersCache(config.Cache{Size: 1}, rep, log)
		if err != nil {
			t.Fatalf("error creating cache: %v", err)
		}

		_, err = cache.Create(context.Background(), testOrder)
		if err == nil {
			t.Fatalf("error was not propagated")
		}

		if !errors.Is(err, decorateeErr) {
			t.Fatal("unexpected error was propagated")
		}
	})

	t.Run("success - eviction on lru size limit achieved", func(t *testing.T) {
		prevOrder := &orders.Order{
			ID: "prev-order",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rep := mocks.NewMockRepository(ctrl)

		rep.EXPECT().
			Create(gomock.Any(), gomock.Eq(prevOrder)).
			Return(prevOrder, nil).
			Times(1)

		rep.EXPECT().
			Create(gomock.Any(), gomock.Eq(testOrder)).
			Return(testOrder, nil).
			Times(1)

		cache, err := NewOrdersCache(config.Cache{Size: 1}, rep, log)
		if err != nil {
			t.Fatalf("error creating cache: %v", err)
		}

		_, err = cache.Create(context.Background(), prevOrder)
		if err != nil {
			t.Fatalf("error while creating first order: %v", err)
		}

		_, err = cache.Create(context.Background(), testOrder)
		if err != nil {
			t.Fatalf("error while creating second order: %v", err)
		}

		if containsPrev := cache.lru.Contains(prevOrder.ID); containsPrev {
			t.Fatalf("value was not evicted")
		}

		if containsNew := cache.lru.Contains(testOrder.ID); !containsNew {
			t.Fatalf("new value was not saved")
		}
	})
}

func TestOrdersCache_Prefill(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.DiscardHandler)

	t.Run("decoratee error is propagated", func(t *testing.T) {
		const size = 5
		decorateeErr := errors.New("some obscure error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rep := mocks.NewMockRepository(ctrl)

		rep.EXPECT().
			ListRecent(gomock.Any(), size).
			Return(nil, decorateeErr).
			Times(1)

		cache, err := NewOrdersCache(config.Cache{Size: size}, rep, log)
		if err != nil {
			t.Fatalf("error creating cache: %v", err)
		}

		err = cache.Prefill(context.Background())
		if err == nil {
			t.Fatal("cache did not return any error")
		}

		if !errors.Is(err, decorateeErr) {
			t.Fatal("returned error does not wrap decoratee error")
		}
	})

	t.Run("caches values on decoratee success", func(t *testing.T) {
		const size = 1

		testOrder := orders.Order{
			ID:        "someid",
			CreatedAt: time.Now(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		rep := mocks.NewMockRepository(ctrl)

		rep.EXPECT().
			ListRecent(gomock.Any(), size).
			Return([]orders.Order{testOrder}, nil).
			Times(1)

		cache, err := NewOrdersCache(config.Cache{Size: size}, rep, log)
		if err != nil {
			t.Fatalf("error creating cache: %v", err)
		}

		if err := cache.Prefill(context.Background()); err != nil {
			t.Fatal("cache errored")
		}

		order, ok := cache.lru.Get(testOrder.ID)
		if !ok {
			t.Fatal("internal cache does not contain order")
		}

		if !reflect.DeepEqual(testOrder, *order) {
			t.Fatal("cache stored unexpected order")
		}
	})
}
