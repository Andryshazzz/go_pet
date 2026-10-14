package dto

import (
	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/entity"
)

// ToOrderEntityItem converts DTO items to entity items.
func ToOrderEntityItems(items []OrderItemDTO) []entity.OrderItem {
	result := make([]entity.OrderItem, len(items))
	for i, item := range items {
		result[i] = entity.OrderItem{
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}
	return result
}

// ToOrderDTO converts entity.Order to OrderDTO.
func ToOrderDTO(order entity.Order) OrderDTO {
	items := make([]OrderItemDTO, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderItemDTO{
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return OrderDTO{
		ID:         order.ID,
		UserID:     order.UserID,
		Status:     string(order.Status),
		Items:      items,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
		UpdatedAt:  order.UpdatedAt,
	}
}

// ToOrderDTOs converts a slice of entity.Order to a slice of OrderDTO.
func ToOrderDTOs(orders []entity.Order) []OrderDTO {
	result := make([]OrderDTO, len(orders))
	for i, order := range orders {
		result[i] = ToOrderDTO(order)
	}
	return result
}
