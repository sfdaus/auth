package domain

import (
	"context"

	"prakarsa-app/transport/request"
)

// User ...
type User struct {
	ID           string `json:"id"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	TokenVersion string `json:"token_version"`
	PublicID     string `json:"public_id"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    int64  `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    int64  `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
	DeletedAt    int64  `json:"deleted_at"`
	IsVerified   bool   `json:"is_verified"`
}

type SignUpUsecase interface {
	SignUp(ctx context.Context, request *request.SignUpReq) (accessToken string, userID string, err error)
	VerifyAccount(ctx context.Context, request *request.VerifyAccountReq) (accessToken string, userID string, err error)
}

type SignInUsecase interface {
	SignIn(ctx context.Context, request *request.SignInReq) (accessToken string, userID string, err error)
}

type ForgotPassword interface {
	ForgotPassword(ctx context.Context, request *request.ForgotPasswordReq) (err error)
	VerifyResetPassword(ctx context.Context, request *request.VerifyResetPasswordReq) (valid bool, err error)
	ResetPassword(ctx context.Context, request *request.ResetPasswordReq) (err error)
}

// UserRepository represent the users repository contract
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	UpdateTokenVersionByID(ctx context.Context, tokenVersion string, id string) error
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (User, error)
	GetByUserID(ctx context.Context, userID string) (User, error)
	UpdatePasswordByUserID(ctx context.Context, userID string, newPassword string) error
	VerifyAccountByUserID(ctx context.Context, userID string) error
}
