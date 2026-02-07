package services

import (
	"ecommerce-api/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenDetails struct {
	AccessToken         string    `json:"access_token"`
	RefreshToken        string    `json:"refresh_token"`
	AccessUUID          uuid.UUID `json:"access_uuid"`
	RefreshUUID         uuid.UUID `json:"refresh_uuid"`
	AccessTokenExpires  int64     `json:"access_token_expires"`
	RefreshTokenExpires int64     `json:"refresh_token_expires"`
}

// Generate Token Pair - Access token & Refresh token
func GenerateTokenPair(user *models.User) (*TokenDetails, error) {
	tokenDetails := &TokenDetails{}
	tokenDetails.AccessTokenExpires = time.Now().Add(time.Minute * 15).Unix() // Access token: 15m
	tokenDetails.AccessUUID = uuid.New()

	tokenDetails.RefreshTokenExpires = time.Now().Add(time.Hour * 24 * 7).Unix() // Refresh token: 7d
	tokenDetails.RefreshUUID = uuid.New()

	var err error
	secret := os.Getenv("JWT_SECRET")

	// 1. Create Access token
	accessTokenClaims := jwt.MapClaims{
		"authorized":  true,
		"access_uuid": tokenDetails.AccessUUID,
		"user_id":     user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"exp":         tokenDetails.AccessTokenExpires,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	tokenDetails.AccessToken, err = accessToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	// 2. Create Refresh token
	refreshTokenClaims := jwt.MapClaims{
		"refresh_uuid": tokenDetails.RefreshUUID,
		"user_id":      user.ID,
		"exp":          tokenDetails.RefreshTokenExpires,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenDetails.RefreshToken, err = refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return tokenDetails, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}
