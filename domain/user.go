package domain

import (
	"context"
	"time"
)

// User ...
type User struct {
	ID           int64     `json:"id"`
	Password     string    `json:"password"`
	Username     string    `json:"username"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	TokenVersion string    `json:"token_version"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	DeletedAt    time.Time `json:"deleted_at"`
}

// UserRepository represent the users repository contract
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (User, error)
}
