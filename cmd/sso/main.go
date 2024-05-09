package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Touchme245/sso_server/internal/app"
	"github.com/Touchme245/sso_server/internal/config"
)

func main() {
	cfg := config.MustLoad()
	// fmt.Println(cfg)
	logger := setupLogger(cfg.Env)
	logger.Info("starting application", slog.String("env", cfg.Env))

	application := app.NewApp(logger, cfg.Grpc.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCServ.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServ.Stop()

	logger.Info("application stopped")

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
