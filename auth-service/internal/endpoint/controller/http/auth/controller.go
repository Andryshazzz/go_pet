package auth

import (
	"context"
	"net/http"

	"github.com/Andryshazzz/go_pet/internal/domain/dto"
	usersservice "github.com/Andryshazzz/go_pet/internal/domain/services/user"
	httprequest "github.com/Andryshazzz/go_pet/pkg/httpserver/request"
	httpresponse "github.com/Andryshazzz/go_pet/pkg/httpserver/response"
	httpserver "github.com/Andryshazzz/go_pet/pkg/httpserver/server"
	"github.com/Andryshazzz/go_pet/pkg/logger"
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
			Handler: h.register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Handler: h.login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/auth/refresh",
			Handler: h.refreshToken,
		},
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register with phone number, password, and full name.
// @Description  Returns user data and JWT token pair (access + refresh).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration request"
// @Success      201  {object}  RegisterResponse "Successfully registered"
// @Failure      400  {object}  httpresponse.ErrorsResponse "Bad request"
// @Failure      409  {object}  httpresponse.ErrorsResponse "Phone already registered"
// @Failure      500  {object}  httpresponse.ErrorsResponse "Internal server error"
// @Router       /auth/register [post]
func (h *UsersHTTPHandler) register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := httpresponse.NewHTTPResponseHandler(log, rw)

	var request dto.RegisterRequest

	if err := httprequest.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate request")
		return
	}

	result, err := h.usersService.Register(ctx, usersservice.RegisterUserRequest{
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
		FullName:    request.FullName,
	})

	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to register user")
		return
	}

	response := dto.RegisterResponse{
		User: dto.UserDTO{
			ID:          result.User.ID,
			PhoneNumber: result.User.PhoneNumber,
			FullName:    result.User.FullName,
		},
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresIn:    result.TokenPair.ExpiresIn,
	}

	responseHandler.JSONResponse(response, http.StatusCreated)
}

// Login godoc
// @Summary      Login user
// @Description  Login with phone number and password.
// @Description  Returns user data and JWT token pair (access + refresh).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login request"
// @Success      200  {object}  LoginResponse "Successfully logged in"
// @Failure      400  {object}  httpresponse.ErrorsResponse "Bad request"
// @Failure      401  {object}  httpresponse.ErrorsResponse "Invalid credentials"
// @Failure      404  {object}  httpresponse.ErrorsResponse "User not found"
// @Failure      500  {object}  httpresponse.ErrorsResponse "Internal server error"
// @Router       /auth/login [post]
func (h *UsersHTTPHandler) login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := httpresponse.NewHTTPResponseHandler(log, rw)

	var request dto.LoginRequest
	if err := httprequest.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate request")
		return
	}

	result, err := h.usersService.Login(ctx, usersservice.LoginUserRequest{
		PhoneNumber: request.PhoneNumber,
		Password:    request.Password,
	})
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to login")
		return
	}

	response := dto.LoginResponse{
		User: dto.UserDTO{
			ID:          result.User.ID,
			PhoneNumber: result.User.PhoneNumber,
			FullName:    result.User.FullName,
		},
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresIn:    result.TokenPair.ExpiresIn,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Use a valid refresh token to obtain a new access and refresh token pair.
// @Description  The old refresh token becomes invalid after this operation.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RefreshTokenRequest true "Refresh token request"
// @Success      200  {object}  RefreshTokenResponse "New token pair"
// @Failure      400  {object}  httpresponse.ErrorsResponse "Bad request"
// @Failure      401  {object}  httpresponse.ErrorsResponse "Invalid or expired refresh token"
// @Failure      500  {object}  httpresponse.ErrorsResponse "Internal server error"
// @Router       /auth/refresh [post]
func (h *UsersHTTPHandler) refreshToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := httpresponse.NewHTTPResponseHandler(log, rw)

	var request dto.RefreshTokenRequest
	if err := httprequest.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to decode and validate request")
		return
	}

	result, err := h.usersService.RefreshToken(ctx, usersservice.RefreshTokenRequest{
		RefreshToken: request.RefreshToken,
	})
	if err != nil {
		responseHandler.ErrorResponse(err, "Failed to refresh token")
		return
	}

	response := dto.RefreshTokenResponse{
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresIn:    result.TokenPair.ExpiresIn,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
