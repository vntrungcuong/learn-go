package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store cung cấp interface cho mọi truy vấn database và quản lý transaction.
// Nó nhúng interface Querier (do SQLC sinh ra trong querier.go).
type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

// SQLStore triển khai Store interface sử dụng pgxpool.
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries // Chứa các phương thức SQLC sinh ra
}

// NewStore khởi tạo một SQLStore instance mới.
// Được gọi từ main.go với tham số là *pgxpool.Pool.
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool), // New nhận vào DBTX (pgxpool thỏa mãn interface này)
	}
}

// ExecTx thực thi một function bên trong một database transaction.
// Đảm bảo tính nguyên tử (Atomicity) cho các nghiệp vụ phức tạp.
func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	// Bắt đầu transaction với cấu hình mặc định (Read Committed)
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	// Tạo Queries mới dựa trên transaction hiện tại
	q := New(tx)
	err = fn(q)
	if err != nil {
		// Rollback nếu có lỗi xảy ra trong closure
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit nếu mọi thứ thành công
	return tx.Commit(ctx)
}
