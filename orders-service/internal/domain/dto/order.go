package dto

import (
	"github.com/google/uuid"
	"time"
)

// OrderDTO represents order data in HTTP responses.
type OrderDTO struct {
	ID         uuid.UUID      `json:"id"`
	UserID     uuid.UUID      `json:"user_id"`
	Status     string         `json:"status"`
	Items      []OrderItemDTO `json:"items"`
	TotalPrice float64        `json:"total_price"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// OrderItemDTO represents a single item in the order.
type OrderItemDTO struct {
	Name     string  `json:"name"     validate:"required,min=1,max=255"`
	Quantity int     `json:"quantity" validate:"required,min=1,max=1000"`
	Price    float64 `json:"price"    validate:"required,min=0"`
}
