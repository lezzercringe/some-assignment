package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"order-persistor/internal/config"
	"order-persistor/internal/orders"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
)

type OrdersConsumer struct {
	consumer *kafka.Consumer
	cfg      config.KafkaConsumer

	ordersRepository orders.Repository
	logger           *slog.Logger

	stopCh chan struct{}
}

// NewOrdersConsumer creates a ready-to-use kafka orders consumer, however at the point of creation no subscription is being done.
// Subscription only starts with an explicit call of Run function.
func NewOrdersConsumer(cfg config.KafkaConsumer, ordersRepository orders.Repository, logger *slog.Logger) (*OrdersConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cfg.Servers,
		"group.id":           cfg.GroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})

	if err != nil {
		return nil, err
	}

	return &OrdersConsumer{
		consumer:         c,
		cfg:              cfg,
		ordersRepository: ordersRepository,
		logger:           logger,
	}, nil
}

func (c *OrdersConsumer) Stop() {
	c.stopCh <- struct{}{}
}

// Run subcribes to the topic and starts processing it, blocking the calling coroutine.
func (c *OrdersConsumer) Run(ctx context.Context) error {
	c.stopCh = make(chan struct{}, 1)
	c.logger.Info("started kafka consumer", "cfg", c.cfg)

	if err := c.consumer.Subscribe(c.cfg.Topic, nil); err != nil {
		return fmt.Errorf("could not subscribe to a topic: %w", err)
	}
	defer c.consumer.Close()

	c.logger.Info("subscribed succesfully")

	for {
		select {
		case <-c.stopCh:
			return nil
		default:
			msg, err := c.consumer.ReadMessage(c.cfg.ReadTimeout)
			if err != nil {
				if err.(kafka.Error).IsTimeout() {
					continue
				}

				return fmt.Errorf("reading message: %w", err)
			}

			c.logger.Debug("consumed order from kafka")

			ctx, cancel := context.WithTimeout(ctx, c.cfg.ProcessTimeout)
			defer cancel()

			if err := c.handleMessage(ctx, msg.Value); err != nil {
				c.logger.ErrorContext(ctx,
					"failure processing order from kafka",
					"err", err,
				)

				continue
			}

			_, err = c.consumer.Commit()
			if err != nil {
				return fmt.Errorf("could not commit message: %w", err)
			}
		}
	}
}

func (c *OrdersConsumer) handleMessage(ctx context.Context, body []byte) error {
	var order orders.Order
	if err := json.Unmarshal(body, &order); err != nil {
		c.logger.ErrorContext(ctx,
			"bad json order from kafka",
			"err", err,
			"message", string(body),
		)

		return nil
	}

	if err := validator.New().StructCtx(ctx, order); err != nil {
		c.logger.ErrorContext(ctx,
			"consumed malformed order from kafka",
			"err", err,
			"message", string(body),
		)

		return nil
	}

	if _, err := c.ordersRepository.Create(ctx, &order); err != nil {
		if errors.Is(err, orders.ErrInternalFailure) {
			return fmt.Errorf("failed creating new order: %w", err)
		}

		c.logger.ErrorContext(
			ctx, "failed creating new order",
			"err", err,
			"message", string(body),
		)
	}

	return nil
}
