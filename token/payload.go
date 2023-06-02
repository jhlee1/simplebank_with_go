package token

import (
	"errors"
	"time"
)

// Payload contains the payload data of the token
type Payload struct {
	ID        string
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issue_at"`
	ExpiredAt time.Time `json:"expire_at"`
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("this token has expired")
)
