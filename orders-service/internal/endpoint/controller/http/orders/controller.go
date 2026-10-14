package orders

import (
	"context"
	"net/http"

	"github.com/Andryshazzz/go_pet/orders-service/internal/domain/dto"
	ordersservice "github.com/Andryshazzz/go_pet/orders-service/internal/domain/services/orders"
	apperrors "github.com/Andryshazzz/go_pet/orders-service/pkg/errors"
	httpresponse "github.com/Andryshazzz/go_pet/orders-service/pkg/httpserver/response"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrdersService defines the interface for order business logic.
type OrdersService interface {
	CreateOrder(ctx context.Context, req ordersservice.CreateOrderRequest) (*ordersservice.CreateOrderResponse, error)
	GetOrder(ctx context.Context, req ordersservice.GetOrderRequest) (*ordersservice.GetOrderResponse, error)
	ListOrders(ctx context.Context, req ordersservice.ListOrdersRequest) (*ordersservice.ListOrdersResponse, error)
}

// OrdersHTTPHandler handles HTTP requests for orders.
type OrdersHTTPHandler struct {
	ordersService OrdersService
}

// NewOrdersHTTPHandler creates a new handler.
func NewOrdersHTTPHandler(ordersService OrdersService) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{ordersService: ordersService}
}

// RegisterRoutes registers all order routes.
//
// All routes require JWT authentication (via middleware).
func (h *OrdersHTTPHandler) RegisterRoutes(router *gin.RouterGroup) {
	orders := router.Group("/orders")
	{
		orders.POST("", h.createOrder)
		orders.GET("/:id", h.getOrder)
		orders.GET("", h.listOrders)
	}
}

// createOrder godoc
// @Summary      Create a new order
// @Description  Creates a new order for the authenticated user
// @Description  Sends a message to Kafka when order is created
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateOrderRequest true "Order data"
// @Success      201 {object} dto.CreateOrderResponse "Order created"
// @Failure      400 {object} httpresponse.ErrorsResponse "Bad request"
// @Failure      401 {object} httpresponse.ErrorsResponse "Unauthorized"
// @Failure      500 {object} httpresponse.ErrorsResponse "Internal server error"
// @Router       /orders [post]
func (h *OrdersHTTPHandler) createOrder(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	userID, ok := getUserID(c)
	if !ok {
		httpresponse.ErrorResponse(c, apperrors.ErrUnauthorized, "failed to get user id")
		return
	}

	var request dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpresponse.ErrorResponse(c, err, "failed to decode and validate request")
		return
	}

	items := make([]ordersservice.OrderItemInput, len(request.Items))
	for i, item := range request.Items {
		items[i] = ordersservice.OrderItemInput{
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	result, err := h.ordersService.CreateOrder(ctx, ordersservice.CreateOrderRequest{
		UserID: userID,
		Items:  items,
	})
	if err != nil {
		log.Error("failed to create order")
		httpresponse.ErrorResponse(c, err, "failed to create order")
		return
	}

	response := dto.CreateOrderResponse{
		Order: dto.ToOrderDTO(result.Order),
	}

	c.JSON(http.StatusCreated, response)
}

// getOrder godoc
// @Summary      Get order by ID
// @Description  Returns order details. User can only access their own orders.
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order ID (UUID)"
// @Success      200 {object} dto.GetOrderResponse "Order details"
// @Failure      401 {object} httpresponse.ErrorsResponse "Unauthorized"
// @Failure      403 {object} httpresponse.ErrorsResponse "Access denied"
// @Failure      404 {object} httpresponse.ErrorsResponse "Order not found"
// @Failure      500 {object} httpresponse.ErrorsResponse "Internal server error"
// @Router       /orders/{id} [get]
func (h *OrdersHTTPHandler) getOrder(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		httpresponse.ErrorResponse(c, apperrors.ErrUnauthorized, "failed to get user id")
		return
	}

	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpresponse.ErrorResponse(c, apperrors.ErrInvalidArgument, "invalid order id format")
		return
	}

	result, err := h.ordersService.GetOrder(ctx, ordersservice.GetOrderRequest{
		OrderID: orderID,
		UserID:  userID,
	})
	if err != nil {
		httpresponse.ErrorResponse(c, err, "failed to get order")
		return
	}

	response := dto.GetOrderResponse{
		Order: dto.ToOrderDTO(result.Order),
	}

	c.JSON(http.StatusOK, response)
}

// listOrders godoc
// @Summary      List user's orders
// @Description  Returns all orders belonging to the authenticated user
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.ListOrdersResponse "List of orders"
// @Failure      401 {object} httpresponse.ErrorsResponse "Unauthorized"
// @Failure      500 {object} httpresponse.ErrorsResponse "Internal server error"
// @Router       /orders [get]
func (h *OrdersHTTPHandler) listOrders(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		httpresponse.ErrorResponse(c, apperrors.ErrUnauthorized, "failed to get user id")
		return
	}

	result, err := h.ordersService.ListOrders(ctx, ordersservice.ListOrdersRequest{
		UserID: userID,
	})
	if err != nil {
		httpresponse.ErrorResponse(c, err, "failed to list orders")
		return
	}
	response := dto.ListOrdersResponse{
		Orders: dto.ToOrderDTOs(result.Orders),
		Total:  result.Total,
	}

	c.JSON(http.StatusOK, response)
}

// getUserID extracts user ID from Gin context (set by JWT middleware).
func getUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	id, ok := userID.(uuid.UUID)
	return id, ok
}
