package dto

// RegisterRequest represents the request for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required,alphanumunicode,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Username string `json:"username" binding:"required,alphanumunicode,min=3,max=32"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the request for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresAt    string  `json:"expires_at"`
	User         UserDTO `json:"user"`
}

// UserDTO represents user data in responses
type UserDTO struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles,omitempty"`
}
