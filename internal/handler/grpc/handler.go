package grpc

import (
	"context"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-auth/internal/service"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthGRPC struct {
	log     *zap.SugaredLogger
	tracer  trace.Tracer
	service service.Auth
	pb.UnimplementedAuthServer
}

func NewAuthGRPC(log *zap.SugaredLogger, tracer trace.Tracer, service service.Auth) *AuthGRPC {
	return &AuthGRPC{log: log, tracer: tracer, service: service}
}

func (a *AuthGRPC) Register(ctx context.Context, input *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, span := a.tracer.Start(ctx, "Register")
	defer span.End()

	id, err := a.service.Register(ctx, input)

	if err != nil {
		a.log.Errorf("Register: %v", err.Error())

		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "Register: %v", err)
	}

	return &pb.RegisterResponse{UserId: id}, nil

}

func (a *AuthGRPC) VerifyUser(ctx context.Context, input *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	ctx, span := a.tracer.Start(ctx, "VerifyUser")
	defer span.End()

	err := a.service.VerifyUser(ctx, input)
	if err != nil {
		a.log.Errorf("VerifyUser: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "VerifyUser: %v", err)
	}

	return &pb.VerifyResponse{}, nil

}

func (a *AuthGRPC) CheckVerify(ctx context.Context, input *pb.CheckVerifyRequest) (*pb.CheckVerifyResponse, error) {
	ctx, span := a.tracer.Start(ctx, "CheckVerify")
	defer span.End()

	err := a.service.CheckVerify(ctx, input)

	if err != nil {
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "CheckVerify: %v", err)
	}

	return &pb.CheckVerifyResponse{}, nil
}

func (a *AuthGRPC) Login(ctx context.Context, input *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, span := a.tracer.Start(ctx, "Login")
	defer span.End()

	token, err := a.service.Login(ctx, input)

	if err != nil {
		a.log.Errorf("Login: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "Login: %v", err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (a *AuthGRPC) GetUserByID(ctx context.Context, input *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	ctx, span := a.tracer.Start(ctx, "GetUserByID")
	defer span.End()

	user, err := a.service.GetByUUID(ctx, input.GetUserId())

	if err != nil {
		a.log.Errorf("GetUser: %v", err.Error())
		span.RecordError(err)
		return nil, err
	}

	return &pb.GetUserResponse{
		UserId:     user.UserID.String(),
		Username:   user.Username,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}

func (a *AuthGRPC) ForgotPassword(ctx context.Context, input *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	ctx, span := a.tracer.Start(ctx, "ForgotPassword")
	defer span.End()

	err := a.service.ForgotPassword(ctx, input)
	if err != nil {
		a.log.Errorf("ForgotPassword: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "ForgotPassword: %v", err)
	}

	return &pb.ForgotPasswordResponse{}, nil
}

func (a *AuthGRPC) ResetPassword(ctx context.Context, input *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	ctx, span := a.tracer.Start(ctx, "ResetPassword")
	defer span.End()

	err := a.service.ResetPassword(ctx, input)
	if err != nil {
		a.log.Errorf("ResetPassword: %v", err.Error())
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "ResetPassword: %v", err)
	}

	return &pb.ResetPasswordResponse{}, nil

}
