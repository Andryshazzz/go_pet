package dto

// CreateOrderResponse represents the response after creating an order.
type CreateOrderResponse struct {
	Order OrderDTO `json:"order"`
}

// GetOrderResponse represents the response for getting an order.
type GetOrderResponse struct {
	Order OrderDTO `json:"order"`
}

type ListOrdersResponse struct {
	Orders []OrderDTO `json:"orders"`
	Total  int        `json:"total"`
}
