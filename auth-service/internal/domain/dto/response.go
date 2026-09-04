package dto

// LoginResponse represents the JSON response after successful login.
type LoginResponse struct {
	User            PrivateUserDTO `json:"user"`
	AccessToken     string         `json:"access_token"`
	RefreshToken    string         `json:"refresh_token"`
	AccessExpiresIn int64          `json:"access_expires_in"`
}

// RefreshTokenResponse represents the JSON response with new token pair.
type RefreshTokenResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessExpiresIn int64  `json:"access_expires_in"`
}

// RegisterResponse represents the JSON response after successful registration.
type RegisterResponse struct {
	User            PrivateUserDTO `json:"user"`
	AccessToken     string         `json:"access_token"`
	RefreshToken    string         `json:"refresh_token"`
	AccessExpiresIn int64          `json:"access_expires_in"`
}
