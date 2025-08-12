package jwt

import "github.com/golang-jwt/jwt"

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}
