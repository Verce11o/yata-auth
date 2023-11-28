package service

import (
	"context"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
)

type Auth interface {
	Register(ctx context.Context, input *pb.RegisterRequest) (int, error)
	Login(ctx context.Context, input *pb.LoginRequest) (string, error)
	Logout(ctx context.Context, input *pb.LogoutRequest) (string, error)
}
