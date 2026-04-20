package auth_transport_http

import (
	"net/http"

	core_http_server "github.com/Andryshazzz/go_pet/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface{}

func NewUsersHTTPHandler(
	usersService UsersService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
