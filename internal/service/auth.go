package service

import (
	"bytes"
	"context"
	"github.com/Verce11o/yata-auth/internal/domain"
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

func (a *AuthService) Register(ctx context.Context, input *pb.RegisterRequest) (string, error) {

	input.Password = a.jwtService.GenerateHashPassword(input.Password)

	userID, err := a.repo.Register(ctx, input)

	if err != nil {
		return "", err
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

	token, err := a.jwtService.GenerateToken(user.UserID.String())

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) GetByUUID(ctx context.Context, userID string) (domain.User, error) {
	return a.repo.GetUserByID(ctx, userID)
}
