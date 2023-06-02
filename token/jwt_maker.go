package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const minSecretKeySize = 32

type JWTPayload struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// jwtMaker is a struct that implements the Maker interface
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

// jwtMaker implements the Maker interface
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := JWTPayload{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "simplebank",
			ID:        uuid.NewString(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// print error message

			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JWTPayload{}, keyFunc)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenUnverifiable) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*JWTPayload)

	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &Payload{ID: payload.ID, Username: payload.Username, IssuedAt: payload.IssuedAt.Time, ExpiredAt: payload.ExpiresAt.Time}, nil
}
