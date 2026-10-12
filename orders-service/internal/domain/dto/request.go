package dto

// CreateOrderRequest represents the JSON body for creating an order.
type CreateOrderRequest struct {
	Items []OrderItemDTO `json:"items" validate:"required,min=1,dive"`
}
