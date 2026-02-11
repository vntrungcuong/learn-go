package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Handle JWT authen middleware
func JwtAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get access token
		tokenString, err := extractTokenString(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": err.Error()})
			return
		}

		// 2. Validate token
		claims, err := validateToken(tokenString, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Invalid or expired token"})
			return
		}

		// 3. Get userId from token
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("x-user-id", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"errors": "Invalid token claims"})
			return
		}

		// 4. Execute next action/request
		c.Next()
	}
}

// Get access token from Header -> Authorization
func extractTokenString(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid token format")
	}

	return parts[1], nil
}

// Validate access token, is valid Signing method
func validateToken(tokenString, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate claims and exp
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Validate success
		return claims, nil
	}

	return nil, errors.New(("Invalid token"))
}
