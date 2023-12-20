package grpc_errors

import (
	"context"
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("not found")
	ErrCodeExpired        = errors.New("code is expired")
	ErrCodeInvalid        = errors.New("code is invalid")
	ErrPasswordMismatch   = errors.New("password mismatch")
	ErrGettingCode        = errors.New("error getting code")
	ErrAlreadyVerified    = errors.New("user already verified")
)

func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrEmailExists):
		return codes.AlreadyExists
	case errors.Is(err, ErrInvalidCredentials):
		return codes.Unauthenticated
	case errors.Is(err, ErrCodeExpired):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrAlreadyVerified):
		return codes.AlreadyExists
	case errors.Is(err, ErrPasswordMismatch):
		return codes.InvalidArgument
	case errors.Is(err, ErrNotFound):
		return codes.NotFound
	case errors.Is(err, redis.Nil):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	}
	return codes.Internal
}
