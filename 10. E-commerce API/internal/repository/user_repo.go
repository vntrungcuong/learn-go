package repository

import (
	"context"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/models"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateUser - Thêm người dùng mới vào database
func CreateUser(ctx context.Context, user *models.User) error {
	// Sử dụng RETURNING để lấy lại ID (UUID) và thời gian tạo do DB sinh ra
	// Giúp tiết kiệm 1 câu lệnh SELECT
	query := `
		INSERT INTO users (email, password_hash, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	// QueryRow thực thi lệnh và map kết quả trả về vào struct user
	err := config.DB.QueryRow(ctx, query,
		user.Email,
		user.Password,
		user.FullName,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Xử lý lỗi trùng Email (Unique Constraint Violation)
		// Mã lỗi 23505 là chuẩn của PostgreSQL cho lỗi duplicate key
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("email already exists")
		}
		// Trả về các lỗi khác (mất kết nối, sai cú pháp...)
		return err
	}

	return nil
}

// GetUserByEmail - Tìm người dùng theo Email (Dùng cho Login)
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := config.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FullName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		// Nếu không tìm thấy dòng nào -> Trả về nil (Không phải lỗi hệ thống)
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
