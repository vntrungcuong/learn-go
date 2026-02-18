package infra

import (
	"context"
	"go-auth-system/internal/auth/domain"
	db "go-auth-system/internal/db/sqlc"
)

type authRepository struct {
	store db.Store
}

// NewAuthRepository khởi tạo repository cho Postgres
func NewAuthRepository(store db.Store) domain.AuthRepository {
	return &authRepository{store: store}
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	return r.store.GetUserByEmail(ctx, email)
}

func (r *authRepository) GetUserByID(ctx context.Context, id int64) (db.User, error) {
	return r.store.GetUserByID(ctx, id)
}

func (r *authRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.store.CreateUser(ctx, arg)
}

func (r *authRepository) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	return r.store.UpdateUser(ctx, arg)
}

func (r *authRepository) UpdatePassword(ctx context.Context, arg db.UpdatePasswordParams) error {
	return r.store.UpdatePassword(ctx, arg)
}
