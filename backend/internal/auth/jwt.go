package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidJWTSecret = errors.New("invalid JWT secret")

func GenerateJWT(
	userID,
	role,
	secret string,
	expirationHours int,
) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrInvalidJWTSecret
	}

	if expirationHours <= 0 {
		return "", errors.New("JWT expiration must be greater than zero")
	}

	now := time.Now()

	expiration := now.Add(
		time.Duration(expirationHours) * time.Hour,
	)

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(secret),
	)
}

func ValidateJWT(
	tokenString,
	secret string,
) (*Claims, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, ErrInvalidJWTSecret
	}

	if strings.TrimSpace(tokenString) == "" {
		return nil, jwt.ErrTokenMalformed
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrTokenSignatureInvalid
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
