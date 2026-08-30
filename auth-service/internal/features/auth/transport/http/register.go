package authtransport

import (
	"net/http"

	usersservice "github.com/Andryshazzz/go_pet/internal/features/auth/service"
	httprequest "github.com/Andryshazzz/go_pet/pkg/httpserver/request"
	httpresponse "github.com/Andryshazzz/go_pet/pkg/httpserver/response"
	logger "github.com/Andryshazzz/go_pet/pkg/logger"
)

// RegisterRequest represents the JSON body for user registration.
type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=15,startswith=+" example:"+79991112233"`
	Password    string `json:"password"     validate:"required,min=8,max=72"              example:"securePassword123"`
	FullName    string `json:"full_name"    validate:"required,min=2,max=100"             example:"Ivan Ivanov"`
}

// RegisterResponse represents the JSON response after successful registration.
type RegisterResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

// UserResponse represents public user data (no password hash).
type UserResponse struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
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
func (h *UsersHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := httpresponse.NewHTTPResponseHandler(log, rw)

	var request RegisterRequest

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

	response := RegisterResponse{
		User: UserResponse{
			PhoneNumber: result.User.PhoneNumber,
			FullName:    result.User.FullName,
		},
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresIn:    result.TokenPair.ExpiresIn,
	}

	responseHandler.JSONResponse(response, http.StatusCreated)
}
