package services

import (
	"context"
	"ecommerce-api/internal/models"
	"ecommerce-api/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
func LoginUser(ctx context.Context, input models.LoginInput) (string, error) {
	// 1. Get user from DB by email
	user, err := repository.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// 2. Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// 3. Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// Hàm nội bộ tạo JWT
func generateJWT(user *models.User) (string, error) {
	claims := models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token sống 1 ngày
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
