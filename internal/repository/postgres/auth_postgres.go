package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Verce11o/yata-auth/internal/domain"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
)

type AuthPostgres struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func NewAuthPostgres(db *sqlx.DB, tracer trace.Tracer) *AuthPostgres {
	return &AuthPostgres{db: db, tracer: tracer}
}

func (s *AuthPostgres) Register(ctx context.Context, input *pb.RegisterRequest) (string, error) {
	ctx, span := s.tracer.Start(ctx, "authPostgres.Register")
	defer span.End()

	var userID uuid.UUID

	q := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING user_id"

	stmt, err := s.db.PreparexContext(ctx, q)

	if err != nil {
		return "", err
	}

	err = stmt.QueryRowxContext(ctx, input.GetUsername(), input.GetEmail(), input.GetPassword()).Scan(&userID)

	var pgErr *pq.Error
	ok := errors.As(err, &pgErr)

	if ok {
		if pgErr.Code == "23505" {
			return "", grpc_errors.ErrEmailExists
		}
	}

	if err != nil {
		return "", err
	}
	return userID.String(), nil
}

func (s *AuthPostgres) AddVerificationCode(ctx context.Context, codeType string, code string, userID string) error {
	ctx, span := s.tracer.Start(ctx, "authPostgres.AddVerificationCode")
	defer span.End()

	q := "INSERT INTO verification_codes (type, code, user_id) VALUES ($1, $2, $3)"

	err := s.db.QueryRowxContext(ctx, q, codeType, code, userID).Err()

	if err != nil {
		return err
	}

	return nil

}

func (s *AuthPostgres) GetVerificationCode(ctx context.Context, codeID string) (*domain.VerificationCode, error) {
	ctx, span := s.tracer.Start(ctx, "authPostgres.GetVerificationCode")
	defer span.End()

	var verificationCode domain.VerificationCode

	q := "SELECT type, code, user_id, expire_date FROM verification_codes WHERE code = $1"

	err := s.db.QueryRowxContext(ctx, q, codeID).StructScan(&verificationCode)

	if err != nil {
		return nil, err
	}

	return &verificationCode, nil
}

func (s *AuthPostgres) ClearVerificationCode(ctx context.Context, userID string, codeType string) error {
	ctx, span := s.tracer.Start(ctx, "authPostgres.ClearVerificationCode")
	defer span.End()

	q := "DELETE FROM verification_codes WHERE user_id = $1 AND type = $2"

	_, err := s.db.ExecContext(ctx, q, userID, codeType)

	if err != nil {
		return err
	}

	return nil
}

func (s *AuthPostgres) VerifyUser(ctx context.Context, userID string) (*domain.User, error) {
	ctx, span := s.tracer.Start(ctx, "authPostgres.VerifyUser")
	defer span.End()

	q := "UPDATE users SET is_verified = true, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1 RETURNING *"

	var user domain.User

	if err := s.db.QueryRowxContext(ctx, q, userID).StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *AuthPostgres) GetUser(ctx context.Context, email string) (domain.User, error) {
	ctx, span := s.tracer.Start(ctx, "authPostgres.GetUser")
	defer span.End()

	var user domain.User

	q := "SELECT * FROM users WHERE email = $1"

	err := s.db.QueryRowxContext(ctx, q, email).StructScan(&user)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, sql.ErrNoRows
	}

	if err != nil {
		return domain.User{}, err
	}

	return user, nil

}

func (s *AuthPostgres) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	ctx, span := s.tracer.Start(ctx, "authPostgres.GetUserByID")
	defer span.End()

	var user domain.User

	q := "SELECT * FROM users WHERE user_id = $1"

	span.AddEvent("main query")
	err := s.db.QueryRowxContext(ctx, q, userID).StructScan(&user)

	span.AddEvent("checking for error")
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, sql.ErrNoRows
	}

	if err != nil {
		return domain.User{}, grpc_errors.ErrInvalidCredentials
	}

	span.AddEvent("returning user")
	return user, nil

}

func (s *AuthPostgres) UpdatePassword(ctx context.Context, userID string, password string) error {
	ctx, span := s.tracer.Start(ctx, "authPostgres.GetUserByID")
	defer span.End()

	q := "UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE user_id = $2"

	res, err := s.db.ExecContext(ctx, q, password, userID)

	if err != nil {
		return nil
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil

}
