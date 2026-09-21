package dto

// LoginRequest represents the JSON body for user login.
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=15,startswith=+" example:"+79991112233"`
	Password    string `json:"password"     validate:"required,min=8,max=72"              example:"securePassword123"`
}

// RefreshTokenRequest represents the JSON body for token refresh.
type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"  validate:"required" example:"eyJhbG..."`
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbG..."`
}

// RegisterRequest represents the JSON body for user registration.
type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,min=10,max=15,startswith=+" example:"+79991112233"`
	Password    string `json:"password"     validate:"required,min=8,max=72"              example:"securePassword123"`
	FullName    string `json:"full_name"    validate:"required,min=2,max=100"             example:"Ivan Ivanov"`
}
