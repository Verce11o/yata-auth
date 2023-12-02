package app

import (
	"fmt"
	"github.com/Verce11o/yata-auth/config"
	authGrpc "github.com/Verce11o/yata-auth/internal/handler/grpc"
	"github.com/Verce11o/yata-auth/internal/lib/auth_jwt"
	"github.com/Verce11o/yata-auth/internal/lib/logger"
	"github.com/Verce11o/yata-auth/internal/repository/postgres"
	"github.com/Verce11o/yata-auth/internal/repository/redis"
	"github.com/Verce11o/yata-auth/internal/service"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	log := logger.NewLogger()
	cfg := config.LoadConfig()

	db := postgres.NewPostgres(cfg)
	repo := postgres.NewAuthPostgres(db)

	rdb := redis.NewRedis(cfg)
	redis := redis.NewAuthRedis(rdb)

	s := grpc.NewServer()

	authService := service.NewAuthService(log, repo, redis, auth_jwt.MakeJWTService(cfg.App.JWT))

	pb.RegisterAuthServer(s, authGrpc.NewAuthGRPC(log, authService))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.App.Port))
	if err != nil {
		log.Info("failed to listen: %v", err)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Infof("error while listen server: %s", err)
		}
	}()

	log.Info(fmt.Sprintf("server listening at %s", lis.Addr().String()))

	defer log.Sync()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.GracefulStop()

	if err := db.Close(); err != nil {
		log.Infof("error while close db: %s", err)
	}

}
