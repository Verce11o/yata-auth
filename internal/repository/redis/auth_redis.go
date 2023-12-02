package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/yata-auth/internal/domain"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	userTTL = 3600
)

type AuthRedis struct {
	client *redis.Client
}

func NewAuthRedis(client *redis.Client) *AuthRedis {
	return &AuthRedis{client: client}
}

func (r *AuthRedis) GetByIdCtx(ctx context.Context, key string) (*domain.User, error) {
	userBytes, err := r.client.Get(ctx, r.createKey(key)).Bytes()

	if err != nil {
		if err != redis.Nil {
			return nil, grpc_errors.ErrNotFound
		}
		return nil, err
	}

	var user *domain.User

	if err = json.Unmarshal(userBytes, &user); err != nil {
		return nil, err
	}

	return user, nil

}

func (r *AuthRedis) SetByIdCtx(ctx context.Context, key string, user *domain.User) error {
	userBytes, err := json.Marshal(user)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createKey(key), userBytes, time.Second*time.Duration(userTTL)).Err()
}

func (r *AuthRedis) DeleteUserCtx(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.createKey(key)).Err()
}

func (r *AuthRedis) createKey(key string) string {
	return fmt.Sprintf("user:%s", key)
}
