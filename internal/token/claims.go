package token

import (
	"fmt"
	"time"

	"github.com/Harman6282/attendance-app/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	ID   string     `json:"id"`
	Role store.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewUserClaims(id string, role store.Role, duration time.Duration) (*JWTClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}

	return &JWTClaims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
