package service

import (
	"bytes"
	"context"
	"github.com/Verce11o/yata-auth/internal/domain"
	"github.com/Verce11o/yata-auth/internal/lib/auth_jwt"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-auth/internal/repository"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AuthService struct {
	log        *zap.SugaredLogger
	tracer     trace.Tracer
	repo       repository.Repository
	redis      repository.RedisRepository
	jwtService auth_jwt.JWTService
}

func NewAuthService(log *zap.SugaredLogger, tracer trace.Tracer, repo repository.Repository, redis repository.RedisRepository, jwtService auth_jwt.JWTService) *AuthService {
	return &AuthService{log: log, tracer: tracer, repo: repo, redis: redis, jwtService: jwtService}
}

func (a *AuthService) Register(ctx context.Context, input *pb.RegisterRequest) (string, error) {
	ctx, span := a.tracer.Start(ctx, "authService.Register")
	defer span.End()

	input.Password = a.jwtService.GenerateHashPassword(input.Password)

	userID, err := a.repo.Register(ctx, input)

	if err != nil {
		return "", err
	}

	return userID, nil
}

func (a *AuthService) Login(ctx context.Context, input *pb.LoginRequest) (string, error) {
	ctx, span := a.tracer.Start(ctx, "authService.Login")
	defer span.End()

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
	ctx, span := a.tracer.Start(ctx, "authService.GetByUUID")
	defer span.End()

	cachedUser, err := a.redis.GetByIdCtx(ctx, userID)

	if err != nil {
		a.log.Errorf("cannot get user by id in redis: %v", err.Error())
	}

	if cachedUser != nil {
		return *cachedUser, nil
	}

	user, err := a.repo.GetUserByID(ctx, userID)

	if err != nil {
		a.log.Errorf("cannot get user by id in postgres: %v", err.Error())
		return domain.User{}, err
	}

	if err := a.redis.SetByIdCtx(ctx, userID, &user); err != nil {
		a.log.Errorf("cannot set user by id in redis: %v", err.Error())
	}

	return user, nil
}
