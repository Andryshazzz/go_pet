package authtransport

import (
	"context"
	"net/http"

	"github.com/Andryshazzz/go_pet/internal/core/domain"
	httpserver "github.com/Andryshazzz/go_pet/internal/core/transport/http/server"
)

// UsersService defines the interface for user business logic operations.
type UsersService interface {
	// CreateUser persists a new user and returns the created user with
	// any server-generated fields populated.
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

// UsersHTTPHandler handles HTTP requests for user-related endpoints.
type UsersHTTPHandler struct {
	usersService UsersService
}

// NewUsersHTTPHandler creates a new UsersHTTPHandler with the given service.
func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

// Routes returns the HTTP route definitions for user endpoints.
// These routes are registered with the API version router in main.
//
// Endpoints:
//   - POST /users — create a new user
func (h *UsersHTTPHandler) Routes() []httpserver.Route {
	return []httpserver.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
