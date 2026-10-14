package ordersservice

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	apperrors "github.com/Andryshazzz/go_pet/orders-service/pkg/errors"
	"github.com/google/uuid"
)

// GetOrderRequest contains the data needed to get an order.
type GetOrderRequest struct {
	OrderID uuid.UUID
	UserID  uuid.UUID // Used for authorization
}

// GetOrderResponse contains the found order.
type GetOrderResponse struct {
	Order entity.Order
}

// GetOrder retrieves an order by ID.
func (s *OrdersService) GetOrder(
	ctx context.Context,
	req GetOrderRequest,
) (*GetOrderResponse, error) {
	order, err := s.ordersRepository.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("get order by id: %w", err)
	}

	if !order.IsOwnedBy(req.UserID) {
		return nil, fmt.Errorf(
			"order %s does not belong to user %s: %w",
			req.OrderID,
			req.UserID,
			apperrors.ErrOrderAccessDenied,
		)
	}

	return &GetOrderResponse{
		Order: order,
	}, nil
}
