package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcDelivery interface {
	HealthChecker
	Register(registry grpc.ServiceRegistrar)
}

type GrpcConfig struct {
	Host string
	Port string
}

type GrpcApp struct {
	config GrpcConfig
	server *grpc.Server
	logger *slog.Logger
}

func InterceptorLogger(logger *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		logger.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func NewGrpcApp(config GrpcConfig, logger *slog.Logger, delivery ...GrpcDelivery) *GrpcApp {
	recoveryOpt := recovery.WithRecoveryHandlerContext(
		func(ctx context.Context, p interface{}) error {
			logger.ErrorContext(ctx, fmt.Sprintf("panic: %s\n\n%s", p, string(debug.Stack())))
			return status.Errorf(codes.Internal, "%s", p)
		},
	)

	serverOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger)),
			recovery.UnaryServerInterceptor(recoveryOpt),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger)),
			recovery.StreamServerInterceptor(recoveryOpt),
		),
	}

	server := grpc.NewServer(serverOpts...)
	reflection.Register(server)

	for _, d := range delivery {
		d.Register(server)
	}

	return &GrpcApp{
		config: config,
		server: server,
		logger: logger,
	}
}

func (app *GrpcApp) Start() error {
	listener, err := net.Listen("tcp", app.config.Host+":"+app.config.Port)
	if err != nil {
		return errors.Wrap(err, "listen tcp")
	}

	return errors.Wrap(app.server.Serve(listener), "start grpc app")
}

func (app *GrpcApp) Shutdown() {
	app.server.GracefulStop()
}
