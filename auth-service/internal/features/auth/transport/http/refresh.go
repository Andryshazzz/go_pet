package authtransport

import (
	"net/http"

	logger "github.com/Andryshazzz/go_pet/internal/core/logger"
	httprequest "github.com/Andryshazzz/go_pet/internal/core/transport/http/request"
	httpresponse "github.com/Andryshazzz/go_pet/internal/core/transport/http/response"
	"github.com/Andryshazzz/go_pet/internal/features/auth/service"
)

// RefreshTokenRequest represents the JSON body for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// RefreshTokenResponse represents the JSON response with new token pair.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
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
func (h *UsersHTTPHandler) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := httpresponse.NewHTTPResponseHandler(log, rw)

	var request RefreshTokenRequest
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

	response := RefreshTokenResponse{
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresIn:    result.TokenPair.ExpiresIn,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}