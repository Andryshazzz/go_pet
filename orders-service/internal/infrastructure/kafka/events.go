package kafkainfra

import (
	"time"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	"github.com/google/uuid"
)

// OrderCreatedEvent is published when a new order is created.
type OrderCreatedEvent struct {
	OrderID    uuid.UUID          `json:"order_id"`
	UserID     uuid.UUID          `json:"user_id"`
	Status     string             `json:"status"`
	Items      []entity.OrderItem `json:"items"`
	TotalPrice float64            `json:"total_price"`
	CreatedAt  time.Time          `json:"created_at"`
}

// OrderStatusChangedEvent is published when an order's status changes.
type OrderStatusChangedEvent struct {
	OrderID   uuid.UUID `json:"order_id"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}