package grpc

import (
	"context"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-auth/internal/service"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type AuthGRPC struct {
	log     *zap.SugaredLogger
	service service.Auth
	pb.UnimplementedAuthServer
}

func NewAuthGRPC(log *zap.SugaredLogger, service service.Auth) *AuthGRPC {
	return &AuthGRPC{log: log, service: service}
}

func (a *AuthGRPC) Register(ctx context.Context, input *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	id, err := a.service.Register(ctx, input)

	if err != nil {
		a.log.Errorf("Register: %v", err.Error())

		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "Register: %v", err)
	}

	return &pb.RegisterResponse{UserId: int64(id)}, nil

}

func (a *AuthGRPC) Login(ctx context.Context, input *pb.LoginRequest) (*pb.LoginResponse, error) {

	token, err := a.service.Login(ctx, input)

	if err != nil {
		a.log.Errorf("Login: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "Login: %v", err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

//
//func (a *AuthGRPC) Logout(ctx context.Context, input *pb.LogoutRequest) (*pb.LogoutResponse, error) {
//
//}
