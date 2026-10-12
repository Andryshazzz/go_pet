package apperrors

type AppError struct {
	HTTPCode int
	Code     string
	Message  string
	Details  map[string]any
}

// Error implements the error interface.
// Returns the human-readable message.
func (e AppError) Error() string {
	return e.Message
}

// Is allows errors.Is to compare AppError by Code.
// This makes AppError work correctly with errors.Is().
func (e AppError) Is(target error) bool {
	t, ok := target.(AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Common application errors used across all layers.
var (
	// ErrNotFoundOrder indicates that the requested order does not exist.
	ErrNotFoundOrder = AppError{
		HTTPCode: 404,
		Code:     "order_not_found",
		Message:  "order not found",
		Details:  nil,
	}

	// ErrOrderAccessDenied indicates that the user does not own the order.
	ErrOrderAccessDenied = AppError{
		HTTPCode: 403,
		Code:     "order_access_denied",
		Message:  "you do not have access to this order",
		Details:  nil,
	}

	// ErrOrderCannotBeModified indicates that the order is in a state
	// that does not allow modification (e.g., already completed).
	ErrOrderCannotBeModified = AppError{
		HTTPCode: 409,
		Code:     "order_cannot_be_modified",
		Message:  "order cannot be modified in its current state",
		Details:  nil,
	}

	// ErrInvalidArgument indicates that the provided input is invalid.
	ErrInvalidArgument = AppError{
		HTTPCode: 422,
		Code:     "invalid_argument",
		Message:  "invalid argument provided",
		Details:  nil,
	}

	// ErrUnauthorized indicates that authentication is required or failed.
	// Used when JWT is missing, invalid, or expired.
	ErrUnauthorized = AppError{
		HTTPCode: 401,
		Code:     "unauthorized",
		Message:  "unauthorized",
		Details:  nil,
	}

	// ErrInvalidToken indicates that the JWT token is invalid, expired, or malformed.
	ErrInvalidToken = AppError{
		HTTPCode: 401,
		Code:     "invalid_token",
		Message:  "invalid or expired token",
		Details:  nil,
	}

	// ErrForbidden indicates that the user is authenticated but does not
	// have permission to perform the action.
	ErrForbidden = AppError{
		HTTPCode: 403,
		Code:     "forbidden",
		Message:  "access forbidden",
		Details:  nil,
	}
)
