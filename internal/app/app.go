package app

import (
	"fmt"
	"github.com/Verce11o/yata-auth/config"
	authGrpc "github.com/Verce11o/yata-auth/internal/handler/grpc"
	"github.com/Verce11o/yata-auth/internal/lib/auth_jwt"
	"github.com/Verce11o/yata-auth/internal/lib/email/rabbitmq"
	"github.com/Verce11o/yata-auth/internal/lib/logger"
	"github.com/Verce11o/yata-auth/internal/lib/metrics/trace"
	"github.com/Verce11o/yata-auth/internal/repository/postgres"
	"github.com/Verce11o/yata-auth/internal/repository/redis"
	"github.com/Verce11o/yata-auth/internal/service"
	pb "github.com/Verce11o/yata-protos/gen/go/sso"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	log := logger.NewLogger()
	cfg := config.LoadConfig()
	tracer := trace.InitTracer("yata-auth")
	//tracer := trace.InitTracer("Yata-Auth")

	// Init storages
	db := postgres.NewPostgres(cfg)
	repo := postgres.NewAuthPostgres(db, tracer.Tracer)

	rdb := redis.NewRedis(cfg)
	redis := redis.NewAuthRedis(rdb, tracer.Tracer)

	// Init rabbitmq client
	amqpConn := rabbitmq.NewAmqpConnection(cfg.RabbitMQ)
	emailPublisher := rabbitmq.NewEmailPublisher(amqpConn, log, tracer.Tracer, cfg.RabbitMQ)

	authService := service.NewAuthService(log, tracer.Tracer, repo, redis, emailPublisher, cfg.App.EmailEndpoint, auth_jwt.MakeJWTService(cfg.App.JWT))

	s := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithTracerProvider(tracer.Provider),
		otelgrpc.WithPropagators(propagation.TraceContext{}),
	)))

	pb.RegisterAuthServer(s, authGrpc.NewAuthGRPC(log, tracer.Tracer, authService))

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
