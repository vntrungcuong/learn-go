package usecase

import (
	"context"
	"errors"
	"time"

	"go-auth-system/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
	jwtSecret      string
}

func NewUserUsecase(u domain.UserRepository, timeout time.Duration, secret string) domain.UserUsecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: timeout,
		jwtSecret:      secret,
	}
}

func (u *userUsecase) Register(c context.Context, email, password, fullname string) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// Check exist
	_, err := u.userRepo.GetByEmail(ctx, email)
	if err == nil {
		// No error <=> Found user data by email
		return errors.New("Email already exists")
	}

	if !errors.Is(err, domain.ErrNotFound) {
		// Other error, return system error
		return err
	}

	// Hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPass),
		FullName: fullname,
	}

	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Login(c context.Context, email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
		"iat":     time.Now().Unix(), // Issued at
	})

	t, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", errors.New("cannot generate token")
	}

	return t, nil
}

func (u *userUsecase) GetProfile(c context.Context, userID string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.userRepo.GetByID(ctx, userID)
}
