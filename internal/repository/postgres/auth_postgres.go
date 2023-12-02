package postgres

import (
	"context"
	"errors"
	"github.com/Verce11o/yata-auth/internal/domain"
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

	stmt, err := s.db.PreparexContext(ctx, q)

	if err != nil {
		return 0, err
	}

	err = stmt.QueryRowxContext(ctx, input.GetUsername(), input.GetEmail(), input.GetPassword()).Scan(&id)

	var pgErr *pq.Error
	ok := errors.As(err, &pgErr)

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

func (s *AuthPostgres) GetUser(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	q := "SELECT * FROM users WHERE email = $1"

	err := s.db.QueryRowxContext(ctx, q, email).StructScan(&user)

	if err != nil {
		return domain.User{}, grpc_errors.ErrInvalidCredentials
	}

	return user, nil

}
