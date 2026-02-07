package services

import (
	"context"
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser - Xử lý logic đăng ký
func RegisterUser(ctx context.Context, input models.RegisterInput) (*models.User, error) {
	// 1. Hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 2. Mapping fields from input to User model
	user := models.User{
		Email:    input.Email,
		Password: string(hashedPwd),
		FullName: input.FullName,
		Role:     "user",
	}

	// 3. Call repository to create user
	err = repository.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// LoginUser - Xử lý logic đăng nhập & cấp token
func LoginUser(ctx context.Context, input models.LoginInput) (*TokenDetails, error) {
	// 1. Get user from DB by email
	user, err := repository.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	// 3. Create access token & refresh token
	tokens, err := GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// 4. Save refresh token to DB
	err = repository.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken, time.Unix(tokens.RefreshTokenExpires, 0))
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// Create Access token from Refresh token
func RefreshToken(ctx context.Context, refreshToken string) (*TokenDetails, error) {
	// 1. Verify format JWT
	token, err := ValidateToken(refreshToken)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 2. Check token is existed in DB or expired
	userID, err := repository.CheckRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("token revoked or expired")
	}

	// 3. Get user info to create token
	user := &models.User{
		ID:   userID,
		Role: "user",
	}

	// 4. Generate new token pair
	newTokens, err := GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// 5. Delete old token and save new token
	_ = repository.RevokeRefreshToken(ctx, refreshToken)
	_ = repository.SaveRefreshToken(ctx, userID, newTokens.RefreshToken, time.Unix(newTokens.RefreshTokenExpires, 0))

	return newTokens, nil
}

// Delete Refresh token
func LogoutUser(ctx context.Context, refreshToken string) error {
	return repository.RevokeRefreshToken(ctx, refreshToken)
}

// Forgot password
func ForgotPassword(ctx context.Context, email string) error {
	// 1. Check user if existed
	user, err := repository.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		// Not show error to reduce leak user info
		return nil
	}

	// 2. Generate reset token
	resetToken := uuid.New().String()

	// 3. Save reset token to DB
	err = repository.SetResetToken(ctx, email, resetToken)
	if err != nil {
		return err
	}

	// 4. Send email to user with reset token
	fmt.Printf("Reset link: http://localhost:3000/reset-password?token=%s\n", resetToken)
	return nil
}

// Reset password
func ResetPassword(ctx context.Context, input models.ResetPasswordInput) error {
	// 1. Verify token
	email, err := repository.VerifyResetToken(ctx, input.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	// 2. Hash new password
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	// 3. Update DB
	return repository.UpdatePassword(ctx, email, string(hashedPwd))
}
