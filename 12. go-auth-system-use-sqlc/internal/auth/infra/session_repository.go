package infra

import (
	"context"
	"time"

	"go-auth-system/internal/auth/domain"
	"github.com/redis/go-redis/v9"
)

type sessionRepository struct {
	redisClient *redis.Client
}

// NewSessionRepository khởi tạo repository cho Redis
func NewSessionRepository(redisClient *redis.Client) domain.SessionRepository {
	return &sessionRepository{redisClient: redisClient}
}

func (r *sessionRepository) SetSession(ctx context.Context, key string, value string, duration time.Duration) error {
	return r.redisClient.Set(ctx, key, value, duration).Err()
}

func (r *sessionRepository) GetSession(ctx context.Context, key string) (string, error) {
	return r.redisClient.Get(ctx, key).Result()
}

func (r *sessionRepository) DelSession(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
