package repository

import (
	"context"
	"ecommerce-api/internal/config"
	"time"

	"github.com/google/uuid"
)

// Save Refresh token to database
func SaveRefreshToken(ctx context.Context, userId uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := config.DB.Exec(ctx, query, userId, token, expiresAt)
	return err
}

// Delete Refresh token (Logout)
func RevokeRefreshToken(ctx context.Context, token string) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE token = $1
	`
	_, err := config.DB.Exec(ctx, query, token)
	return err
}

// Check Refresh token is existed and not expired
func CheckRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	query := `
		SELECT user_id
		FROM refresh_tokens
		WHERE token = $1
		AND expires_at > NOW()
	`
	var userId uuid.UUID
	err := config.DB.QueryRow(ctx, query, token).Scan(&userId)
	return userId, err
}

// Save Reset token to database
func SetResetToken(ctx context.Context, email string, token string) error {
	// Reset token expires in 15 minutes
	expiresAt := time.Now().Add(time.Minute * 15)
	query := `
		UPDATE users
		SET 
			reset_token = $1, 
			reset_token_exp = $2
		WHERE email = $3
	`
	_, err := config.DB.Exec(ctx, query, token, expiresAt, email)
	return err
}

// Verify Reset token
func VerifyResetToken(ctx context.Context, token string) (string, error) {
	query := `
		SELECT email
		FROM users
		WHERE reset_token = $1
		AND reset_token_exp > NOW()
	`
	var email string
	err := config.DB.QueryRow(ctx, query, token).Scan(&email)
	return email, err
}

// Update password
func UpdatePassword(ctx context.Context, email string, hashPwd string) error {
	query := `
		UPDATE user
		SET
			password_hash = $1,
			reset_token = NULL,
			reset_token_exp = NULL
		WHERE email = $2
	`
	_, err := config.DB.Exec(ctx, query, hashPwd, email)
	return err
}
