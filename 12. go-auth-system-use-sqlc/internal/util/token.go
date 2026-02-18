// Token utils
package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToen  = errors.New("token has expired")
)

type Payload struct {
	UserID               int64 `json:"user_id"`
	jwt.RegisteredClaims `json:"claims"`
}

func CreateToken(userID int64, duration time.Duration, secretKey string) (string, error) {
	payload := &Payload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string, secretKey string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrInvalidToken
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToen
		}
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
