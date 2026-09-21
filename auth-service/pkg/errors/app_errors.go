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

// Common application errors used across all layers.
var (
	// ErrNotFoundUser indicates that the requested user does not exist.
	ErrNotFoundUser = AppError{
		HTTPCode: 404,
		Code:     "user_not_found",
		Message:  "user not found",
		Details:  nil,
	}

	// ErrInvalidArgument indicates that the provided input is invalid.
	ErrInvalidArgument = AppError{
		HTTPCode: 422,
		Code:     "invalid_argument",
		Message:  "invalid argument provided",
		Details:  nil,
	}

	// ErrUserAlreadyExists indicates that a user with the same phone number
	// already exists in the system.
	ErrUserAlreadyExists = AppError{
		HTTPCode: 409,
		Code:     "user_already_exists",
		Message:  "user with this phone number already exists",
		Details:  nil,
	}

	// ErrInvalidCredentials indicates that the provided phone number or password
	// does not match any user in the system.
	ErrInvalidCredentials = AppError{
		HTTPCode: 401,
		Code:     "invalid_credentials",
		Message:  "invalid phone number or password",
		Details:  nil,
	}

	// ErrInvalidToken indicates that the JWT token is invalid, expired, or malformed.
	ErrInvalidToken = AppError{
		HTTPCode: 401,
		Code:     "invalid_token",
		Message:  "invalid or expired token",
		Details:  nil,
	}

	// ErrTokenMismatch indicates that access and refresh tokens belong
	// to different users.
	ErrTokenMismatch = AppError{
		HTTPCode: 401,
		Code:     "token_mismatch",
		Message:  "access and refresh token mismatch",
		Details:  nil,
	}
)
