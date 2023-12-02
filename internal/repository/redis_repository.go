package repository

import (
	"context"
	"github.com/Verce11o/yata-auth/internal/domain"
)

type RedisRepository interface {
	GetByIdCtx(ctx context.Context, key string) (*domain.User, error)
	SetByIdCtx(ctx context.Context, key string, user *domain.User) error
	DeleteUserCtx(ctx context.Context, key string) error
}
