package domain

import (
	"context"
	"time"

	db "go-auth-system/internal/db/sqlc"
)

// AuthRepository định nghĩa các thao tác với Database
type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByID(ctx context.Context, id int64) (db.User, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error)
	UpdatePassword(ctx context.Context, arg db.UpdatePasswordParams) error
}

// SessionRepository định nghĩa các thao tác với Redis
type SessionRepository interface {
	SetSession(ctx context.Context, key string, value string, duration time.Duration) error
	GetSession(ctx context.Context, key string) (string, error)
	DelSession(ctx context.Context, key string) error
}
