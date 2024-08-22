package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	SessionID string `json:"sessionId"`
	jwt.RegisteredClaims
}

func newClaims(sessionID, userID, firstName, lastName string, duration time.Duration) *Claims {
	claims := Claims{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			// ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return &claims
}
