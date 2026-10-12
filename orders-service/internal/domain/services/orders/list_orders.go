package ordersservice

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	"github.com/google/uuid"
)

// ListOrdersRequest contains the data needed to list orders.
type ListOrdersRequest struct {
	UserID uuid.UUID
}

// ListOrdersResponse contains the list of orders.
type ListOrdersResponse struct {
	Orders []entity.Order
	Total  int
}

// ListOrders retrieves all orders for a user.
func (s *OrdersService) ListOrders(
	ctx context.Context,
	req ListOrdersRequest,
) (*ListOrdersResponse, error) {
	orders, err := s.ordersRepository.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("get orders by user id: %w", err)
	}

	return &ListOrdersResponse{
		Orders: orders,
		Total:  len(orders),
	}, nil
}
