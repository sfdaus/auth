package jwt

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"prakarsa-app/config/constant"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: secretKey,
		issuer:    "auth",
	}
}

func (s *jwtService) GenerateToken(ctx context.Context, userID string, tokenVersion string) (token string, err error) {
	claims := &jwtCustomClaims{
		userID,
		tokenVersion,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = t.SignedString([]byte(s.secretKey))
	return
}

func (s *jwtService) GenerateForgotPasswordToken(ctx context.Context, userID string) (token string, jti string, err error) {
	jti = uuid.NewString()
	claims := &forgotPasswordJWTCustomClaims{
		userID,
		jwt.StandardClaims{
			Id:        jti,
			ExpiresAt: time.Now().Add(constant.FORGOT_PASSWORD_EXPIRES).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
			Audience:  "forgot_password",
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = t.SignedString([]byte(s.secretKey))
	return
}

func (s *jwtService) ValidateToken(ctx context.Context, tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}

func (s *jwtService) ParseForgotPasswordToken(ctx context.Context, tokenStr string) (*forgotPasswordJWTCustomClaims, error) {
	// Parse + verify signature
	tkn, err := jwt.ParseWithClaims(tokenStr, &forgotPasswordJWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// pastikan algoritmanya sesuai (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %T", token.Method)
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tkn.Claims.(*forgotPasswordJWTCustomClaims)
	if !ok || !tkn.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Validasi issuer & audience (defense-in-depth)
	if claims.Issuer != s.issuer {
		return nil, fmt.Errorf("invalid issuer")
	}
	if !claims.VerifyAudience("forgot_password", true) {
		return nil, fmt.Errorf("invalid audience")
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
