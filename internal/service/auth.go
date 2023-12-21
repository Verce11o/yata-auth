package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/yata-auth/internal/domain"
	"github.com/Verce11o/yata-auth/internal/lib/auth_jwt"
	"github.com/Verce11o/yata-auth/internal/lib/email"
	"github.com/Verce11o/yata-auth/internal/lib/grpc_errors"
	"github.com/Verce11o/yata-auth/internal/repository"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

const (
	EmailCodeType = "email"
	PassCodeType  = "password"
)

type AuthService struct {
	log                   *zap.SugaredLogger
	tracer                trace.Tracer
	repo                  repository.Repository
	redis                 repository.RedisRepository
	emailPublisher        email.EmailPublisher
	emailEndpoint         string
	passwordResetEndpoint string
	jwtService            auth_jwt.JWTService
}

func NewAuthService(log *zap.SugaredLogger, tracer trace.Tracer, repo repository.Repository, redis repository.RedisRepository, emailPublisher email.EmailPublisher, emailEndpoint string, passwordResetEndpoint string, jwtService auth_jwt.JWTService) *AuthService {
	return &AuthService{log: log, tracer: tracer, repo: repo, redis: redis, emailPublisher: emailPublisher, emailEndpoint: emailEndpoint, passwordResetEndpoint: passwordResetEndpoint, jwtService: jwtService}
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

// todo эта функуия = send email, нужна еще get code в которой будут проверки
func (a *AuthService) VerifyUser(ctx context.Context, input *pb.VerifyRequest) error {
	ctx, span := a.tracer.Start(ctx, "authService.VerifyUser")
	defer span.End()

	user, err := a.GetByUUID(ctx, input.GetUserId())

	if err != nil {
		return err
	}

	if user.IsVerified {
		return grpc_errors.ErrAlreadyVerified
	}

	code := uuid.NewString()

	err = a.repo.AddVerificationCode(ctx, EmailCodeType, code, input.GetUserId())

	if err != nil {
		return err
	}

	SendEmailRequest := domain.SendUserEmailRequest{
		Type: EmailCodeType,
		To:   user.Email,
		Code: fmt.Sprintf("%v?code=%v", a.emailEndpoint, code),
	}

	messageBytes, err := json.Marshal(SendEmailRequest)

	if err != nil {
		return err
	}

	err = a.emailPublisher.Publish(ctx, messageBytes)

	if err != nil {
		return err
	}

	return nil

}

func (a *AuthService) CheckVerify(ctx context.Context, input *pb.CheckVerifyRequest) error {
	ctx, span := a.tracer.Start(ctx, "authService.CheckVerify")
	defer span.End()

	code, err := a.repo.GetVerificationCode(ctx, input.GetCode())

	if err != nil || code == nil {
		a.log.Infof("cannot get verified code by id in postgres: %v", err.Error())
		return grpc_errors.ErrGettingCode
	}

	if code.Code.String() != input.GetCode() {
		a.log.Infof("invalid code: %v and %v", code.Code.String(), input.GetCode())
		return grpc_errors.ErrCodeInvalid
	}

	if time.Now().UTC().After(code.ExpireDate) {
		a.log.Errorf("code is expired")
		return grpc_errors.ErrCodeExpired
	}

	user, err := a.repo.VerifyUser(ctx, code.UserID.String())

	if err != nil {
		a.log.Errorf("cannot verify user: %v", err.Error())
		return err
	}

	err = a.repo.ClearVerificationCode(ctx, code.UserID.String(), code.Type)

	if err != nil {
		a.log.Errorf("cannot clear user codes")
		return err
	}

	err = a.redis.DeleteUserCtx(ctx, user.UserID.String()) // maybe use set instead

	if err != nil {
		a.log.Errorf("cannot delete user in redis")
		return err
	}

	return nil

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

func (a *AuthService) ForgotPassword(ctx context.Context, input *pb.ForgotPasswordRequest) error {
	ctx, span := a.tracer.Start(ctx, "authService.ForgotPassword")
	defer span.End()

	user, err := a.GetByUUID(ctx, input.GetUserId())

	if err != nil {
		return err
	}

	code := uuid.NewString()

	err = a.repo.AddVerificationCode(ctx, PassCodeType, code, input.GetUserId())

	if err != nil {
		return err
	}

	SendEmailRequest := domain.SendUserEmailRequest{
		Type: PassCodeType,
		To:   user.Email,
		Code: fmt.Sprintf("%v?code=%v", a.passwordResetEndpoint, code),
	}

	messageBytes, err := json.Marshal(SendEmailRequest)

	if err != nil {
		return err
	}

	err = a.emailPublisher.Publish(ctx, messageBytes)

	if err != nil {
		return err
	}

	return nil

}

func (a *AuthService) VerifyPassword(ctx context.Context, input *pb.VerifyPasswordRequest) error {
	ctx, span := a.tracer.Start(ctx, "authService.VerifyPassword")
	defer span.End()

	code, err := a.repo.GetVerificationCode(ctx, input.GetCode())

	if err != nil || code == nil {
		a.log.Infof("cannot get pass code by id in postgres: %v", err.Error())
		return grpc_errors.ErrGettingCode
	}

	if code.Code.String() != input.GetCode() {
		a.log.Infof("invalid code: %v and %v", code.Code.String(), input.GetCode())
		return grpc_errors.ErrCodeInvalid
	}

	if time.Now().UTC().After(code.ExpireDate) {
		a.log.Errorf("code is expired")
		return grpc_errors.ErrCodeExpired
	}

	return nil
}
func (a *AuthService) ResetPassword(ctx context.Context, input *pb.ResetPasswordRequest) error {
	ctx, span := a.tracer.Start(ctx, "authService.ResetPassword")
	defer span.End()

	if input.GetPassword() != input.GetPasswordRe() {
		a.log.Errorf("password mismatch")
		return grpc_errors.ErrPasswordMismatch
	}

	hashedPass := a.jwtService.GenerateHashPassword(input.GetPassword())

	err := a.repo.UpdatePassword(ctx, input.GetUserId(), hashedPass)

	if err != nil {
		a.log.Errorf("cannot update user password: %v", err.Error())
		return err
	}

	err = a.repo.ClearVerificationCode(ctx, input.GetCode(), PassCodeType)

	if err != nil {
		a.log.Errorf("cannot clear user codes")
		return err
	}

	return nil

}
