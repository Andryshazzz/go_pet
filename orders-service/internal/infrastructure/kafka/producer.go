package kafkainfra

import (
	"context"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/kafka"
)

// Producer is a domain-specific Kafka producer for order events.
type Producer struct {
	client             *kafka.Producer
	ordersCreatedTopic string
	orderStatusTopic   string
}

// Config holds domain-specific Kafka topics.
type Config struct {
	OrdersCreatedTopic string
	OrderStatusTopic   string
}

// NewProducer creates a new order events producer.
func NewProducer(client *kafka.Producer, config Config) *Producer {
	return &Producer{
		client:             client,
		ordersCreatedTopic: config.OrdersCreatedTopic,
		orderStatusTopic:   config.OrderStatusTopic,
	}
}

// SendOrderCreated publishes an "order created" event.
func (p *Producer) SendOrderCreated(ctx context.Context, order entity.Order) error {
	event := OrderCreatedEvent{
		OrderID:    order.ID,
		UserID:     order.UserID,
		Status:     string(order.Status),
		Items:      order.Items,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
	}

	return p.client.PublishJSON(ctx, p.ordersCreatedTopic, order.ID.String(), event)
}

// SendOrderStatusChanged publishes an "order status changed" event.
func (p *Producer) SendOrderStatusChanged(ctx context.Context, order entity.Order) error {
	event := OrderStatusChangedEvent{
		OrderID:   order.ID,
		UserID:    order.UserID,
		Status:    string(order.Status),
		UpdatedAt: order.UpdatedAt,
	}

	return p.client.PublishJSON(ctx, p.orderStatusTopic, order.ID.String(), event)
}
