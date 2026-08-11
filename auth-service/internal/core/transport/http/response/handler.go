package httpresponse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apperrors "github.com/Andryshazzz/go_pet/internal/core/errors"
	logger "github.com/Andryshazzz/go_pet/internal/core/logger"
	"go.uber.org/zap"
)

// HTTPResponseHandler provides a convenient way to write HTTP responses
// with proper error mapping, logging, and JSON encoding.

type HTTPResponseHandler struct {
	log *logger.Logger
	rw  http.ResponseWriter
}

// NewHTTPResponseHandler creates a new HTTPResponseHandler.
func NewHTTPResponseHandler(log *logger.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

// JSONResponse writes a JSON response with the given status code.
func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

// ErrorResponse maps an application error to an HTTP status code,
// logs it at the appropriate level, and writes an error response.
//
// Error mapping:
//   - ErrInvalidArgument → 400 (Warn log level)
//   - ErrNotFound        → 404 (Debug log level)
//   - ErrConflict        → 409 (Warn log level)
//   - Other errors       → 500 (Error log level)
func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, apperrors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, apperrors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, apperrors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

// PanicResponse handles recovered panics by logging them and returning
// a 500 Internal Server Error response.
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("Unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

// errorResponse writes an ErrorsResponse as JSON with the given status code.
func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := ErrorsResponse{
		Error:   err.Error(),
		Message: msg,
	}

	h.JSONResponse(
		response,
		statusCode,
	)
}
