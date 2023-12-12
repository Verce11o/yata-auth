package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/yata-auth/internal/domain"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const (
	userTTL = 3600
)

type AuthRedis struct {
	client *redis.Client
	tracer trace.Tracer
}

func NewAuthRedis(client *redis.Client, tracer trace.Tracer) *AuthRedis {
	return &AuthRedis{client: client, tracer: tracer}
}

func (r *AuthRedis) GetByIdCtx(ctx context.Context, key string) (*domain.User, error) {
	ctx, span := r.tracer.Start(ctx, "authRedis.GetByIdCtx")
	defer span.End()

	span.AddEvent("getting bytes")
	userBytes, err := r.client.Get(ctx, r.createKey(key)).Bytes()

	if err != nil {
		if err != redis.Nil {
			return nil, grpc_errors.ErrNotFound
		}
		return nil, err
	}

	var user *domain.User

	span.AddEvent("unmarshal")
	if err = json.Unmarshal(userBytes, &user); err != nil {
		return nil, err
	}

	return user, nil

}

func (r *AuthRedis) SetByIdCtx(ctx context.Context, key string, user *domain.User) error {
	ctx, span := r.tracer.Start(ctx, "authRedis.SetByIdCtx")
	defer span.End()

	span.AddEvent("marshal")
	userBytes, err := json.Marshal(user)

	if err != nil {
		return err
	}

	span.AddEvent("set in db")
	return r.client.Set(ctx, r.createKey(key), userBytes, time.Second*time.Duration(userTTL)).Err()
}

func (r *AuthRedis) DeleteUserCtx(ctx context.Context, key string) error {
	ctx, span := r.tracer.Start(ctx, "authRedis.DeleteUserCtx")
	defer span.End()

	return r.client.Del(ctx, r.createKey(key)).Err()
}

func (r *AuthRedis) createKey(key string) string {
	return fmt.Sprintf("user:%s", key)
}
