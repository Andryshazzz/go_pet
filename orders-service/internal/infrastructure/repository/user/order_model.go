package ordersrepository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	utils "github.com/Andryshazzz/go_pet/orders-service/pkg/database/postgres/pgxutil"
	apperrors "github.com/Andryshazzz/go_pet/orders-service/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// OrderModel represents the orders table row in the database.
//
// Fields:
//   - ID: unique identifier (UUID)
//   - UserID: owner of the order
//   - Status: current order status (pending, processing, completed, cancelled)
//   - TotalPrice: total price of the order
//   - CreatedAt: creation timestamp
//   - UpdatedAt: last modification timestamp
type OrderModel struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Status     string
	TotalPrice float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// toEntity converts the database model to a domain entity.
func (m OrderModel) toEntity() entity.Order {
	return entity.Order{
		ID:         m.ID,
		UserID:     m.UserID,
		Status:     entity.OrderStatus(m.Status),
		TotalPrice: m.TotalPrice,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// scanOrder reads an order from a row or rows scanner.
// Works with both pgx.Row and pgx.Rows.
//
// Returns ErrNotFoundOrder if no rows were returned.
func scanOrder(s utils.Scanner) (entity.Order, error) {
	var model OrderModel

	err := s.Scan(
		&model.ID,
		&model.UserID,
		&model.Status,
		&model.TotalPrice,
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{}, apperrors.ErrNotFoundOrder
		}
		return entity.Order{}, fmt.Errorf("scan order: %w", err)
	}

	return model.toEntity(), nil
}

// scanOrderItem reads a single order item from a scanner.
func scanOrderItem(s utils.Scanner) (entity.OrderItem, error) {
	var item entity.OrderItem
	err := s.Scan(
		&item.Name,
		&item.Quantity,
		&item.Price,
	)
	if err != nil {
		return entity.OrderItem{}, fmt.Errorf("scan order item: %w", err)
	}
	return item, nil
}
