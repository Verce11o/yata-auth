package service

import (
	"context"
	"github.com/Verce11o/yata-auth/config"
	"github.com/Verce11o/yata-auth/internal/lib/auth_jwt"
	"github.com/Verce11o/yata-auth/internal/repository"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"go.uber.org/zap"
)

type AuthService struct {
	log        *zap.SugaredLogger
	repo       *repository.Repository
	jwtService *auth_jwt.JWTService
}

func NewAuthService(log *zap.SugaredLogger, repo *repository.Repository, config *config.Config) *AuthService {
	return &AuthService{log: log, repo: repo, jwtService: auth_jwt.NewJWTService(config.App.JWT)}

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

	userID, err := a.repo.Login(ctx, input)

	if err != nil {
		return "", err
	}

	token, err := a.jwtService.GenerateToken(userID)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) Logout(ctx context.Context, input *pb.LogoutRequest) (string, error) {
	//TODO implement me
	panic("implement me")
}
