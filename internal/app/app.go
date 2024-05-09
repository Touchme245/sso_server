package app

import (
	"log/slog"
	"time"

	grpcApp "github.com/Touchme245/sso_server/internal/app/grpc"
	"github.com/Touchme245/sso_server/internal/services/auth"
	"github.com/Touchme245/sso_server/internal/storage/sqlite"
)

type App struct {
	GRPCServ *grpcApp.App
}

func NewApp(logger *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.NewStorage(storagePath)

	if err != nil {
		panic(err)
	}
	autService := auth.NewAuthService(logger, storage, storage, storage, tokenTTL)
	grpcapp := grpcApp.NewApp(logger, grpcPort, autService)
	return &App{
		GRPCServ: grpcapp,
	}
}
