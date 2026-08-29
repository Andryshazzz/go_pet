package authtransport

import (
	"context"
	"net/http"

	usersservice "github.com/Andryshazzz/go_pet/internal/features/auth/service"
	httpserver "github.com/Andryshazzz/go_pet/pkg/httpserver/server"
)

// UsersService defines the interface for user business logic operations.
type UsersService interface {
	Register(ctx context.Context, req usersservice.RegisterUserRequest) (*usersservice.RegisterUserResponse, error)
	Login(ctx context.Context, req usersservice.LoginUserRequest) (*usersservice.LoginUserResponse, error)
	RefreshToken(ctx context.Context, req usersservice.RefreshTokenRequest) (*usersservice.RefreshTokenResponse, error)
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

// PublicRoutes returns routes that don't require authentication.
func (h *UsersHTTPHandler) Routes() []httpserver.Route {
	return []httpserver.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.RefreshToken,
		},
	}
}
