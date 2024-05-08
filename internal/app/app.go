package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/Touchme245/sso_server/internal/app/grpc"
)

type App struct {
	GRPCServ *grpcApp.App
}

func NewApp(logger *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	grpcapp := grpcApp.NewApp(logger, grpcPort)
	return &App{
		GRPCServ: grpcapp,
	}
}
