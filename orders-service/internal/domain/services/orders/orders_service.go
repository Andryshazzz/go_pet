package ordersservice

import (
	"context"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
	"github.com/google/uuid"
)

// OrdersRepository defines the interface for order persistence operations.
// It is defined here (in the service layer) following the Dependency
// Inversion Principle — the service declares what it needs.
type OrdersRepository interface {
	// Create persists a new order and returns it with server-generated fields.
	Create(ctx context.Context, order entity.Order) (entity.Order, error)

	// GetByID retrieves an order by ID.
	// Returns ErrNotFoundOrder if the order does not exist.
	GetByID(ctx context.Context, id uuid.UUID) (entity.Order, error)

	// GetByUserID retrieves all orders for a user.
	// Returns an empty slice if the user has no orders.
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)

	// UpdateStatus updates the status of an order.
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.OrderStatus) error
}

// KafkaProducer defines the interface for sending Kafka messages.
// Used to notify other services about order events.
type KafkaProducer interface {
	// SendOrderCreated sends an order creation event.
	SendOrderCreated(ctx context.Context, order entity.Order) error

	// SendOrderStatusChanged sends a status change event.
	SendOrderStatusChanged(ctx context.Context, order entity.Order) error
}

// OrdersService contains business logic for orders.
type OrdersService struct {
	ordersRepository OrdersRepository
	kafkaProducer    KafkaProducer
}

// NewOrdersService creates a new OrdersService with the given dependencies.
func NewOrdersService(
	ordersRepository OrdersRepository,
	kafkaProducer KafkaProducer,
) *OrdersService {
	return &OrdersService{
		ordersRepository: ordersRepository,
		kafkaProducer:    kafkaProducer,
	}
}