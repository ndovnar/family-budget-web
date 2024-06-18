package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
)

func (auth *Auth) createToken(sessionID, userID, firstName, lastName string, duration time.Duration) (string, error) {
	claims := newClaims(sessionID, userID, firstName, lastName, duration)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtToken.SignedString([]byte(auth.secretKey))
	token, err := jwtToken.SignedString([]byte(auth.secretKey))
	return token, err
}

func (auth *Auth) CreateAccessToken(sessionID, userID, firstName, lastName string) (string, error) {
	return auth.createToken(sessionID, userID, firstName, lastName, auth.accessTokenDuration)
}

func (auth *Auth) CreateRefreshToken(sessionID, userID, firstName, lastName string) (string, error) {
	return auth.createToken(sessionID, userID, firstName, lastName, auth.refreshTokenDuration)
}

func (auth *Auth) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(auth.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return &Claims{}, ErrInvalidToken
	}

	return claims, nil
}
