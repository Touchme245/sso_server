package main

import (
	"log/slog"
	"os"

	"github.com/Touchme245/sso_server/internal/app"
	"github.com/Touchme245/sso_server/internal/config"
)

func main() {
	cfg := config.MustLoad()
	// fmt.Println(cfg)
	logger := setupLogger(cfg.Env)
	logger.Info("starting application", slog.String("env", cfg.Env))

	application := app.NewApp(logger, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)
	application.GRPCServ.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log

}
