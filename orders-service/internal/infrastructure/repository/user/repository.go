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
func (r *OrdersRepository) Create(
	ctx context.Context,
	order entity.Order,
) (entity.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return entity.Order{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	orderQuery := `
		INSERT INTO orders (id, user_id, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = tx.Exec(ctx, orderQuery,
		order.ID,
		order.UserID,
		string(order.Status),
		order.TotalPrice,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return entity.Order{}, fmt.Errorf("insert order: %w", err)
	}

	itemQuery := `
		INSERT INTO order_items (order_id, name, quantity, price)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQuery,
			order.ID,
			item.Name,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return entity.Order{}, fmt.Errorf("insert order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.Order{}, fmt.Errorf("commit transaction: %w", err)
	}

	return order, nil
}

// GetByID retrieves an order by its ID.
func (r *OrdersRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Order, error) {
	query := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	order, err := scanOrder(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		return entity.Order{}, err
	}

	items, err := r.getOrderItems(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}

	order.Items = items
	return order, nil
}

// GetByUserID retrieves all orders for a user.
func (r *OrdersRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Order, error) {
	query := `
		SELECT id, user_id, status, total_price, created_at, updated_at
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
	orderIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		order, err := scanOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
		orderIDs = append(orderIDs, order.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	if len(orders) == 0 {
		return orders, nil
	}

	itemsByOrder, err := r.getOrderItemsBatch(ctx, orderIDs)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		orders[i].Items = itemsByOrder[orders[i].ID]
	}

	return orders, nil
}

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

// getOrderItemsBatch fetches items for multiple orders in one query.
func (r *OrdersRepository) getOrderItemsBatch(
	ctx context.Context,
	orderIDs []uuid.UUID,
) (map[uuid.UUID][]entity.OrderItem, error) {
	if len(orderIDs) == 0 {
		return map[uuid.UUID][]entity.OrderItem{}, nil
	}

	query := `
		SELECT order_id, name, quantity, price
		FROM order_items
		WHERE order_id = ANY($1)
	`

	rows, err := r.pool.Query(ctx, query, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]entity.OrderItem)
	for rows.Next() {
		var orderID uuid.UUID
		var item entity.OrderItem

		if err := rows.Scan(&orderID, &item.Name, &item.Quantity, &item.Price); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}

		result[orderID] = append(result[orderID], item)
	}

	return result, rows.Err()
}

// getOrderItems is a helper that fetches all items for an order.
func (r *OrdersRepository) getOrderItems(
	ctx context.Context,
	orderID uuid.UUID,
) ([]entity.OrderItem, error) {
	query := `
		SELECT name, quantity, price
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query order items: %w", err)
	}
	defer rows.Close()

	items := make([]entity.OrderItem, 0)
	for rows.Next() {
		item, err := scanOrderItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order items: %w", err)
	}

	return items, nil
}
