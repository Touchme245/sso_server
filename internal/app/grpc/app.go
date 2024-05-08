package grpcApp

import (
	"fmt"
	"log/slog"
	"net"

	// authgrpc "sso/internal/grpc/auth"

	authgrpc "github.com/Touchme245/sso_server/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(log *slog.Logger, port int) *App {
	grpcServer := grpc.NewServer()
	authgrpc.Register(grpcServer)
	return &App{
		logger:     log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

func (a *App) Run() error {
	const op = "grpcApp.Run"
	log := a.logger.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf("%d", a.port))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcApp.Stop"
	a.logger.With(slog.String("op", op)).Info("stopping grpc service", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}

func (a *App) MustRun() {
	err := a.Run()
	if err != nil {
		panic(err)
	}
}
