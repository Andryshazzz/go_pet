package authtransport

import (
	"net/http"

	httprequest "github.com/Andryshazzz/go_pet/pkg/transport/http/request"
	"github.com/Andryshazzz/go_pet/internal/features/auth/service"
	logger "github.com/Andryshazzz/go_pet/pkg/logger"
	httpresponse "github.com/Andryshazzz/go_pet/pkg/transport/http/response"
)

// LoginRequest represents the JSON body for user login.
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=15,startswith=+" example:"+79991112233"`
	Password    string `json:"password"     validate:"required,min=8,max=72"              example:"securePassword123"`
}

// LoginResponse represents the JSON response after successful login.
type LoginResponse struct {
	User         usersservice.UserResponse `json:"user"`
	AccessToken  string                    `json:"access_token"`
	RefreshToken string                    `json:"refresh_token"`
	ExpiresIn    int64                     `json:"expires_in"`
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
func (h *UsersHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := httpresponse.NewHTTPResponseHandler(log, rw)

	var request LoginRequest
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

	response := LoginResponse{
		User:         result.User,
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresIn:    result.TokenPair.ExpiresIn,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
