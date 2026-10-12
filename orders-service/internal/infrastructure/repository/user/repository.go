package ordersrepository

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	postgrespool "github.com/Andryshazzz/go_pet/orders-service/pkg/database/postgres/pool"
	apperrors "github.com/Andryshazzz/go_pet/orders-service/pkg/errors"
	"github.com/google/uuid"
)

// OrdersRepository implements data access operations for orders
// using PostgreSQL as the storage backend.
type OrdersRepository struct {
	pool postgrespool.Pool
}

// NewOrdersRepository creates a new OrdersRepository with the given pool.
func NewOrdersRepository(pool postgrespool.Pool) *OrdersRepository {
	return &OrdersRepository{pool: pool}
}

// Create inserts a new order into the database and returns the persisted order.
//
// Implements: ordersservice.OrdersRepository
func (r *OrdersRepository) Create(
	ctx context.Context,
	order entity.Order,
) (entity.Order, error) {
	query := `
		INSERT INTO orders (id, user_id, status, items, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, status, items, total_price, created_at, updated_at;
	`

	row := r.pool.QueryRow(ctx, query,
		order.ID,
		order.UserID,
		string(order.Status),
		order.Items,
		order.TotalPrice,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return scanOrder(row)
}

// GetByID retrieves an order by its ID.
//
// Returns ErrNotFoundOrder if the order does not exist.
func (r *OrdersRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Order, error) {
	query := `
		SELECT id, user_id, status, items, total_price, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	return scanOrder(row)
}

// GetByUserID retrieves all orders for a user.
// Returns an empty slice if the user has no orders.
func (r *OrdersRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Order, error) {
	query := `
		SELECT id, user_id, status, items, total_price, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	orders := make([]entity.Order, 0)

	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	return orders, nil
}

// UpdateStatus updates the status of an order and its updated_at timestamp.
//
// Returns ErrNotFoundOrder if the order does not exist.
func (r *OrdersRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status entity.OrderStatus,
) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	tag, err := r.pool.Exec(ctx, query, string(status), id)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("order %s: %w", id, apperrors.ErrNotFoundOrder)
	}

	return nil
}
