package authtransport

import (
	"github.com/google/uuid"
)

// UserDTO represents public user data in HTTP responses.
// It is independent from service and domain layers.
type UserDTO struct {
	ID          uuid.UUID `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	FullName    string    `json:"full_name"`
}
