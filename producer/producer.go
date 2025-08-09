package producer

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	client *kafka.Producer
	cfg    Config
}

func NewProducer(cfg Config) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Servers,
	})

	if err != nil {
		return nil, err
	}

	return &Producer{
		client: producer,
		cfg:    cfg,
	}, nil
}

func (p *Producer) Start() {
	for e := range p.client.Events() {
		if kafkaErr, ok := e.(kafka.Error); ok {
			fmt.Printf("Error delivering message: %v\n", kafkaErr)
		}
	}
}

func (p *Producer) Produce(msg []byte) error {
	return p.client.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &p.cfg.Topic,
		},
		Value:     msg,
		Timestamp: time.Now(),
	}, nil)
}

func (p *Producer) Close() {
	p.client.Flush(1500)
	p.client.Close()
}
