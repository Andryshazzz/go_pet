package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Config holds Kafka producer configuration.
type Config struct {
	Brokers  []string
	ClientID string
}

// Producer wraps a Confluent Kafka producer.
type Producer struct {
	producer *ckafka.Producer
}

// NewProducer creates a new Kafka producer.
func NewProducer(config Config) (*Producer, error) {
	producer, err := ckafka.NewProducer(&ckafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Brokers, ","),
		"client.id":         config.ClientID,
		"acks":              "all",
	})
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &Producer{producer: producer}, nil
}

// Publish sends a message to the given topic.
func (p *Producer) Publish(
	ctx context.Context,
	topic string,
	key string,
	value []byte,
) error {
	deliveryChan := make(chan ckafka.Event, 1)

	err := p.producer.Produce(&ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: value,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("produce message: %w", err)
	}

	select {
	case e := <-deliveryChan:
		m := e.(*ckafka.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// PublishJSON is a convenience method that marshals the payload to JSON.
func (p *Producer) PublishJSON(
	ctx context.Context,
	topic string,
	key string,
	payload any,
) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	return p.Publish(ctx, topic, key, data)
}

// Close flushes pending messages and closes the producer.
func (p *Producer) Close() {
	p.producer.Flush(5000)
	p.producer.Close()
}
