// AuthMiddleware is a middleware function that validates JWT tokens.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-auth-system/internal/util"
)

func AuthMiddleware(config util.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get Authorization from Request Headers
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			// 2. Check format Bearer <token>
			fields := strings.Fields(authHeader)
			if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// 3. Extract token from fields
			accessToken := fields[1]

			// 4. Validate token
			payload, err := util.VerifyToken(accessToken, config.JWTSecret)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// 5. Set user ID in context
			ctx := context.WithValue(r.Context(), util.UserIDKey, payload.UserID)
			r = r.WithContext(ctx)

			// 6. Call next handler
			next.ServeHTTP(w, r)
		})
	}
}
