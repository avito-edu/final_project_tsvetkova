package jwtutil

import (
	"fmt"
	"swim_service/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64       `json:"user_id"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte
}

func NewTokenService(secret []byte) *TokenService {
	return &TokenService{
		secret: secret,
	}
}

func (ts *TokenService) GenerateAccessToken(userID int64, role domain.Role) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ts.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
