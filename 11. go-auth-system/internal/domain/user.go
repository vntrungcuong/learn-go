package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User entity
type User struct {
	// type:uuid để GORM biết map sang cột UUID của Postgres
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository interface (Repository pattern) -> Interface segregation
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

// UserUseCase interface (Service layer)
type UserUsecase interface {
	Register(ctx context.Context, email, password, fullname string) error
	Login(ctx context.Context, email, password string) (string, error)
	GetProfile(ctx context.Context, userID string) (*User, error)
}
