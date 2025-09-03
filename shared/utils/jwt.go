package utils

import (
	"errors"
	"log"
	"time"

	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/golang-jwt/jwt/v4"
	// "gitlab.com/truemeds-dev-team/truemeds-dev-doctor/truemeds-dev-service/doctorportal-auth-service/shared/constants"
)

func GenerateJwtToken[T any](tokenType constants.TokenType, payload T, expTimeInMinutes int, secret string) (string, error) {

	claims := jwt.MapClaims{
		"payload": payload, // <‑‑ goes into the token exactly as supplied
		"type":    tokenType,
		"exp": time.Now().
			Add(time.Minute * time.Duration(expTimeInMinutes)).
			Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJwtToken(tokenStr string, secret string) (jwt.MapClaims, error) {
	// secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Panic("JWT SECRET is not set in .env file")
	}

	// Parse token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure signing method is expected
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrInvalidSigning
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, constants.ErrExpiredToken
		}
		return nil, constants.ErrInvalidToken
	}

	// Extract and assert claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, constants.ErrInvalidToken

	}

	return claims, nil
}
