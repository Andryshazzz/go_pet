// orders-service/internal/domain/entity/order.go
package entity

import (
	"fmt"
	"time"

	apperrors "github.com/Andryshazzz/go_pet/orders-service/pkg/errors"
	"github.com/google/uuid"
)

// OrderStatus represents the current state of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusCompleted  OrderStatus = "completed"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// OrderItem represents a single item in the order.
type OrderItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

// Order represents a customer's order.
//
// Fields:
//   - ID: unique identifier (UUID)
//   - UserID: owner of the order (from JWT)
//   - Status: current order status
//   - Items: list of items in the order
//   - TotalPrice: total price (calculated from items)
//   - CreatedAt: when order was created
//   - UpdatedAt: when order was last modified
type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Status     OrderStatus
	Items      []OrderItem
	TotalPrice float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewOrder creates a new Order with generated ID and pending status.
func NewOrder(userID uuid.UUID, items []OrderItem) (Order, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Order{}, err
	}

	now := time.Now().UTC()

	return Order{
		ID:         id,
		UserID:     userID,
		Status:     OrderStatusPending,
		Items:      items,
		TotalPrice: calculateTotal(items),
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

// Validate checks that the order is valid.
// Returns ErrInvalidArgument if any field is invalid.
func (o *Order) Validate() error {
	if len(o.Items) == 0 {
		return fmt.Errorf("order must contain at least one item: %w", apperrors.ErrInvalidArgument)
	}

	for i, item := range o.Items {
		if item.Name == "" {
			return fmt.Errorf("item[%d] name cannot be empty: %w", i, apperrors.ErrInvalidArgument)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf(
				"item[%d] quantity must be positive, got %d: %w",
				i, item.Quantity, apperrors.ErrInvalidArgument,
			)
		}
		if item.Price < 0 {
			return fmt.Errorf(
				"item[%d] price cannot be negative, got %f: %w",
				i, item.Price, apperrors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

// IsOwnedBy checks if the order belongs to the given user.
// Used for authorization: user can only see their own orders.
func (o *Order) IsOwnedBy(userID uuid.UUID) bool {
	return o.UserID == userID
}

func calculateTotal(items []OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
