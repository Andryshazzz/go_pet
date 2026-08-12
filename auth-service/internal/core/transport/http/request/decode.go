package httprequest

import (
	"encoding/json"
	"fmt"
	"net/http"

	apperrors "github.com/Andryshazzz/go_pet/internal/core/errors"
	"github.com/go-playground/validator"
)

// requestValidator is a singleton validator instance used for all requests.
var requestValidator = validator.New()

// DecodeAndValidateRequest decodes a JSON request body into the destination
// struct and validates it using the `validate` struct tags.
//
// Returns ErrInvalidArgument if decoding or validation fails.
func DecodeAndValidateRequest(r *http.Request, dest any) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("Decode json: %v: %w", err, apperrors.ErrInvalidArgument)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf("Request validation: %v: %w", err, apperrors.ErrInvalidArgument)
	}

	return nil
}
