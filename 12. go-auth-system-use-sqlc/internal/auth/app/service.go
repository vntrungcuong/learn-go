package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go-auth-system/internal/auth/domain"
	db "go-auth-system/internal/db/sqlc"
	"go-auth-system/internal/util"

	"github.com/google/uuid"
)

// AuthService định nghĩa các nghiệp vụ (Business Use Cases) của hệ thống
type AuthService interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	Register(ctx context.Context, email, password, fullname string) (db.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetProfile(ctx context.Context, userID int64) (db.User, error)
	UpdateProfile(ctx context.Context, userID int64, fullname string) (db.User, error)
	ForgotPassword(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, email, token, newPassword string) error
	Logout(ctx context.Context, userID int64) error
}

type authService struct {
	repo        domain.AuthRepository
	sessionRepo domain.SessionRepository
	config      util.Config
}

// NewAuthService khởi tạo Service với Dependency Injection
func NewAuthService(repo domain.AuthRepository, sessionRepo domain.SessionRepository, config util.Config) AuthService {
	return &authService{
		repo:        repo,
		sessionRepo: sessionRepo,
		config:      config,
	}
}

// 1. Logic Đăng nhập & Tạo Session (Stateful)
func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("user not found or invalid credentials")
	}

	if err := util.CheckPassword(password, user.HashedPassword); err != nil {
		return "", "", errors.New("incorrect password")
	}

	accessToken, _ := util.CreateToken(user.ID, s.config.AccessTokenDuration, s.config.JWTSecret)
	refreshToken, _ := util.CreateToken(user.ID, s.config.RefreshTokenDuration, s.config.JWTSecret)

	// Lưu Refresh Token vào Redis để tối ưu Performance thay vì DB
	redisKey := fmt.Sprintf("refresh_token:%d", user.ID)
	err = s.sessionRepo.SetSession(ctx, redisKey, refreshToken, s.config.RefreshTokenDuration)

	return accessToken, refreshToken, err
}

// 2. Logic Đăng ký & Hash mật khẩu
func (s *authService) Register(ctx context.Context, email, password, fullname string) (db.User, error) {
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	arg := db.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPassword,
		Fullname:       fullname,
	}

	return s.repo.CreateUser(ctx, arg)
}

// 3. Logic Làm mới Token (Refresh Token Rotation)
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	payload, err := util.VerifyToken(refreshToken, s.config.JWTSecret)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Kiểm tra tính hợp lệ của token trong Redis (Stateful check)
	redisKey := fmt.Sprintf("refresh_token:%d", payload.UserID)
	storedToken, err := s.sessionRepo.GetSession(ctx, redisKey)
	if err != nil || storedToken != refreshToken {
		return "", errors.New("token expired or revoked")
	}

	// Chỉ tạo mới Access Token (Stateless)
	return util.CreateToken(payload.UserID, s.config.AccessTokenDuration, s.config.JWTSecret)
}

// 4. Lấy thông tin Profile
func (s *authService) GetProfile(ctx context.Context, userID int64) (db.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

// 5. Cập nhật Profile
func (s *authService) UpdateProfile(ctx context.Context, userID int64, fullname string) (db.User, error) {
	arg := db.UpdateUserParams{
		ID:       userID,
		Fullname: fullname,
	}
	return s.repo.UpdateUser(ctx, arg)
}

// 6. Logic Quên mật khẩu & Tạo Reset Token (Redis TTL)
func (s *authService) ForgotPassword(ctx context.Context, email string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Tránh leak thông tin user (User Enumeration Protection)
		return "", nil
	}

	resetToken := uuid.New().String()
	redisKey := fmt.Sprintf("reset_password:%s", user.Email)

	// Lưu vào Redis với TTL 15 phút
	err = s.sessionRepo.SetSession(ctx, redisKey, resetToken, 15*time.Minute)

	return resetToken, err
}

// 7. Logic Đặt lại mật khẩu & Xóa Token một lần (Atomic)
func (s *authService) ResetPassword(ctx context.Context, email, token, newPassword string) error {
	redisKey := fmt.Sprintf("reset_password:%s", email)
	storedToken, err := s.sessionRepo.GetSession(ctx, redisKey)
	if err != nil || storedToken != token {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, _ := util.HashPassword(newPassword)
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:             user.ID,
		HashedPassword: hashedPassword,
	})

	if err == nil {
		// Xóa token ngay sau khi dùng (One-time use)
		_ = s.sessionRepo.DelSession(ctx, redisKey)
	}

	return err
}

// 8. Logic Đăng xuất (Revoke session)
func (s *authService) Logout(ctx context.Context, userID int64) error {
	redisKey := fmt.Sprintf("refresh_token:%d", userID)
	return s.sessionRepo.DelSession(ctx, redisKey)
}
