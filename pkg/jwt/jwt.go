package jwt

import (
	"errors"
	"strings"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwtv5.RegisteredClaims
}

var (
	ErrEmptySecret   = errors.New("jwt secret must not be empty")
	ErrInvalidExpire = errors.New("jwt expire must be greater than zero")
)

func GenerateToken(secret string, userID uint, role string, expireSeconds int) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrEmptySecret
	}
	if expireSeconds <= 0 {
		return "", ErrInvalidExpire
	}

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Duration(expireSeconds) * time.Second)),
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (*Claims, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, ErrEmptySecret
	}

	token, err := jwtv5.ParseWithClaims(tokenString, &Claims{}, func(token *jwtv5.Token) (any, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, jwtv5.ErrTokenSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwtv5.ErrTokenInvalidClaims
	}

	return claims, nil
}
