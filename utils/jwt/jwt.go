package jwt

import (
	"context"
	"errors"
	"fmt"
	"prakarsa-app/config"
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

func (s *jwtService) ValidateToken(ctx context.Context, tokenString string) (token *jwt.Token, err error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}

func GenerateShortJWT(userID string) (string, error) {
	secret := []byte(config.LoadConfig().JWTSecretKey)
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Duration(config.LoadConfig().AccountVerificationTtl)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateShortJWT(tokenStr string) (userID string, err error) {
	secret := []byte(config.LoadConfig().JWTSecretKey)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid sub")
	}
	return sub, nil
}
