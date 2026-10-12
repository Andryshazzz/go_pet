package ordersservice

import (
	"context"
	"fmt"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateOrderRequest contains the data needed to create an order.
type CreateOrderRequest struct {
	UserID uuid.UUID
	Items  []OrderItemInput
}

// OrderItemInput represents a single item from the request.
type OrderItemInput struct {
	Name     string
	Quantity int
	Price    float64
}

// CreateOrderResponse contains the created order.
type CreateOrderResponse struct {
	Order entity.Order
}

// CreateOrder creates a new order and sends a Kafka event.
func (s *OrdersService) CreateOrder(
	ctx context.Context,
	req CreateOrderRequest,
) (*CreateOrderResponse, error) {
	log := logger.FromContext(ctx)
	items := make([]entity.OrderItem, len(req.Items))
	
	for i, item := range req.Items {
		items[i] = entity.OrderItem{
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	order, err := entity.NewOrder(req.UserID, items)
	if err != nil {
		return nil, fmt.Errorf("create order entity: %w", err)
	}

	if err := order.Validate(); err != nil {
		return nil, fmt.Errorf("validate order: %w", err)
	}

	createdOrder, err := s.ordersRepository.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	if err := s.kafkaProducer.SendOrderCreated(ctx, createdOrder); err != nil {
		log.Warn("failed to send kafka event",
			zap.Error(err),
			zap.String("order_id", createdOrder.ID.String()),
			zap.String("event", "order.created"),
		)
	}

	return &CreateOrderResponse{
		Order: createdOrder,
	}, nil
}
