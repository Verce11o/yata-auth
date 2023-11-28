package repository

import (
	"context"
	"github.com/Verce11o/yata-auth/internal/repository/postgres"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	Register(ctx context.Context, input *pb.RegisterRequest) (int, error)
	Login(ctx context.Context, input *pb.LoginRequest) (int, error)
}

type Repository struct {
	Auth
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth: postgres.NewAuthPostgres(db),
	}
}
