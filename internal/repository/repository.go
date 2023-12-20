package repository

import (
	"context"
	"github.com/Verce11o/yata-auth/internal/domain"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
)

type Repository interface { // maybe refactor to smaller interface?
	Register(ctx context.Context, input *pb.RegisterRequest) (string, error)
	GetUser(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, userID string) (domain.User, error)
	UpdatePassword(ctx context.Context, userID string, password string) error
	AddVerificationCode(ctx context.Context, codeType string, code string, userID string) error
	GetVerificationCode(ctx context.Context, codeID string) (*domain.VerificationCode, error)
	ClearVerificationCode(ctx context.Context, userID string, codeType string) error
	VerifyUser(ctx context.Context, userID string) (*domain.User, error)
}
