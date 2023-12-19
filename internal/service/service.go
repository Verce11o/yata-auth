package service

import (
	"context"
	"github.com/Verce11o/yata-auth/internal/domain"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
)

type Auth interface {
	Register(ctx context.Context, input *pb.RegisterRequest) (string, error)
	VerifyUser(ctx context.Context, input *pb.VerifyRequest) error
	CheckVerify(ctx context.Context, input *pb.CheckVerifyRequest) error
	ForgotPassword(ctx context.Context, input *pb.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, input *pb.ResetPasswordRequest) error
	Login(ctx context.Context, input *pb.LoginRequest) (string, error)
	GetByUUID(ctx context.Context, userID string) (domain.User, error)
}
