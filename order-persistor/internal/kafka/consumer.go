package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"order-persistor/internal/config"
	"order-persistor/internal/orders"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-playground/validator/v10"
)

type Client interface {
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error
	Commit() ([]kafka.TopicPartition, error)
	Close() error
}

type OrdersConsumer struct {
	client Client
	cfg    config.KafkaConsumer

	ordersRepository orders.Repository
	logger           *slog.Logger

	shutdownOnce sync.Once
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
		client:           c,
		cfg:              cfg,
		ordersRepository: ordersRepository,
		logger:           logger,
	}, nil
}

// Run subcribes to the topic and starts processing it, blocking the calling coroutine.
func (c *OrdersConsumer) Run(ctx context.Context) error {
	c.logger.Info("started kafka consumer", "cfg", c.cfg)

	if err := c.client.Subscribe(c.cfg.Topic, nil); err != nil {
		return fmt.Errorf("could not subscribe to a topic: %w", err)
	}
	defer c.closeConsumer()

	// this is needed to unblock main goroutine from being stuck in ReadMessage() call
	// in case of context cancellation
	go func() {
		<-ctx.Done()
		c.closeConsumer()
	}()

	c.logger.Info("subscribed succesfully")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.client.ReadMessage(c.cfg.ReadTimeout)
			if err != nil {
				if err.(kafka.Error).IsTimeout() {
					continue
				}

				c.logger.Error(
					"error reading messages from kafka",
					"err", err, "will start retrying each", c.cfg.ReadFailureBackoff.String(),
				)

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(c.cfg.ReadFailureBackoff):
					continue
				}
			}

			c.logger.Debug("consumed order from kafka")

			// process message inside a closure to be able to use defer cancel()
			err = func() error {
				ctx, cancel := context.WithTimeout(ctx, c.cfg.ProcessTimeout)
				defer cancel()
				return c.handleMessage(ctx, msg.Value)
			}()

			// do not commit in case of internal error (if order was valid)
			if err != nil && !errors.Is(err, errMalformedOrder) {
				c.logger.ErrorContext(ctx,
					"failure processing order from kafka",
					"err", err,
				)

				return err
			}

			_, err = c.client.Commit()
			if err != nil {
				return fmt.Errorf("could not commit message: %w", err)
			}

			continue
		}
	}
}

var errMalformedOrder = errors.New("message was malformed")

// handleMessage processes JSON-encoded order message
// In case of bad json/order returns errMalformedOrder
func (c *OrdersConsumer) handleMessage(ctx context.Context, body []byte) error {
	var order orders.Order
	if err := json.Unmarshal(body, &order); err != nil {
		c.logger.ErrorContext(ctx,
			"bad json order from kafka",
			"err", err,
			"message", string(body),
		)

		return errors.Join(errMalformedOrder, err)
	}

	if err := validator.New().StructCtx(ctx, order); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		c.logger.ErrorContext(ctx,
			"consumed malformed order from kafka",
			"err", err,
			"message", string(body),
		)

		return errors.Join(errMalformedOrder, err)
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

		if ctx.Err() != nil {
			return ctx.Err()
		}

		return errors.Join(errMalformedOrder, err)
	}

	return nil
}

// closeConsumer wraps closing inner consumer into sync.Once
func (c *OrdersConsumer) closeConsumer() {
	c.shutdownOnce.Do(func() {
		c.client.Close()
	})
}
