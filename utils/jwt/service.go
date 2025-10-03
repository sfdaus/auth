package jwt

import (
	"context"

	"github.com/golang-jwt/jwt"
)

type JWTService interface {
	GenerateToken(ctx context.Context, userID string, tokenVersion string) (token string, err error)
	GenerateForgotPasswordToken(ctx context.Context, userID string) (token string, jti string, err error)
	ValidateToken(ctx context.Context, tokenString string) (token *jwt.Token, err error)
	ParseForgotPasswordToken(ctx context.Context, tokenStr string) (*forgotPasswordJWTCustomClaims, error)
}
