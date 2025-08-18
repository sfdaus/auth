package domain

import "context"

// AuthToken ...
type AuthToken struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	UserAgent    string `json:"user_agent"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	IssuedAt     int64  `json:"issued_at"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    int64  `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    int64  `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
	DeletedAt    int64  `json:"deleted_at"`
}

type AuthTokenRepository interface {
	Create(ctx context.Context, authToken *AuthToken) error
}
