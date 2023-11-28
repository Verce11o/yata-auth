package postgres

import (
	"context"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (s *AuthPostgres) Register(ctx context.Context, input *pb.RegisterRequest) (int, error) {
	var id int

	q := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	err := s.db.QueryRowxContext(ctx, q, input.Username, input.Email, input.Password).Scan(&id)

	pgErr, ok := err.(*pq.Error)
	if ok {
		if pgErr.Code == "23505" {
			return 0, grpc_errors.ErrEmailExists
		}
	}

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthPostgres) Login(ctx context.Context, input *pb.LoginRequest) (int, error) {
	var userID int

	q := "SELECT id FROM users WHERE email = $1 OR username = $2 AND password = $3"

	err := s.db.QueryRowxContext(ctx, q, input.GetEmail(), input.GetUsername(), input.GetPassword()).Scan(&userID)

	if err != nil {
		return 0, grpc_errors.ErrInvalidCredentials
	}

	return userID, nil

}
