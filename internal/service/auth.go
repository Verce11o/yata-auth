package service

import (
	"bytes"
	"context"
	"github.com/Verce11o/yata-auth/internal/lib/auth_jwt"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-auth/internal/repository"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"go.uber.org/zap"
)

type AuthService struct {
	log        *zap.SugaredLogger
	repo       repository.Repository
	redis      repository.RedisRepository
	jwtService auth_jwt.JWTService
}

func NewAuthService(log *zap.SugaredLogger, repo repository.Repository, redis repository.RedisRepository, jwtService auth_jwt.JWTService) *AuthService {
	return &AuthService{log: log, repo: repo, redis: redis, jwtService: jwtService}
}

func (a *AuthService) Register(ctx context.Context, input *pb.RegisterRequest) (int, error) {

	input.Password = a.jwtService.GenerateHashPassword(input.Password)

	userID, err := a.repo.Register(ctx, input)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (a *AuthService) Login(ctx context.Context, input *pb.LoginRequest) (string, error) {

	input.Password = a.jwtService.GenerateHashPassword(input.GetPassword())

	user, err := a.repo.GetUser(ctx, input.GetEmail())

	if err != nil {
		return "", err
	}

	if !bytes.Equal([]byte(input.GetPassword()), user.PasswordHash) {
		return "", grpc_errors.ErrInvalidCredentials
	}

	token, err := a.jwtService.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}
