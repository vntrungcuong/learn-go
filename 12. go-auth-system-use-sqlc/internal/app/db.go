package app

import (
	"context"
	"log"
	"time"

	"go-auth-system/internal/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// ConnectPostgres khởi tạo Connection Pool cho PostgreSQL
func ConnectPostgres(config util.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse connection string từ config
	poolConfig, err := pgxpool.ParseConfig(config.PostgreSQLConn)
	if err != nil {
		return nil, err
	}

	// Tối ưu hóa Pool cho Production (Scalability)
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Kiểm tra kết nối
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("✅ Postgres connected via pgxpool")
	return pool, nil
}

// ConnectRedis khởi tạo Redis Client
func ConnectRedis(config util.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	// Kiểm tra kết nối Redis (Fail-fast)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}

	log.Println("✅ Redis connected")
	return client
}
