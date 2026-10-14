package httpresponse

import (
	"errors"
	"net/http"

	apperrors "github.com/Andryshazzz/go_pet/orders-service/pkg/errors"
	"github.com/Andryshazzz/go_pet/orders-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorsResponse represents a standardized error response body.
//
// Example JSON:
//
//	{
//	    "error": "invalid FullName length 1 (must be 2-100): invalid argument",
//	    "message": "failed to create user"
//	}
type ErrorsResponse struct {
	// Error contains the technical error details for debugging.
	Error string `json:"error"`

	// Message contains the human-readable error description.
	Message string `json:"message"`
}

// ErrorResponse maps an application error to an HTTP status code,
// logs it at the appropriate level, and writes an error response.
func ErrorResponse(c *gin.Context, err error, msg string) {
	log := logger.FromContext(c.Request.Context())

	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, apperrors.ErrInvalidArgument):
		statusCode = http.StatusUnprocessableEntity
		logFunc = log.Warn

	case errors.Is(err, apperrors.ErrUnauthorized),
		errors.Is(err, apperrors.ErrInvalidToken):
		statusCode = http.StatusUnauthorized
		logFunc = log.Warn

	case errors.Is(err, apperrors.ErrForbidden),
		errors.Is(err, apperrors.ErrOrderAccessDenied):
		statusCode = http.StatusForbidden
		logFunc = log.Warn

	case errors.Is(err, apperrors.ErrNotFoundOrder):
		statusCode = http.StatusNotFound
		logFunc = log.Debug

	case errors.Is(err, apperrors.ErrOrderCannotBeModified):
		statusCode = http.StatusConflict
		logFunc = log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = log.Error
	}

	logFunc(msg, zap.Error(err))

	c.JSON(statusCode, ErrorsResponse{
		Error:   err.Error(),
		Message: msg,
	})
}
