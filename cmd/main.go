package main

import (
	"context"
	"fmt"
	"link-base/internal/cache"
	"link-base/internal/config"
	"link-base/internal/http"
	"link-base/internal/repository"
	"link-base/internal/server"
	"link-base/internal/service"
	"link-base/pkg/auth"
	"link-base/pkg/database"
	"link-base/pkg/hash"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title LinkBase
// @version 1.0
// @description LinkBase API

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoad()

	fmt.Println("Config: ", cfg)

	logger := setupLogger()

	postgresClient, err := database.NewPostgresClient(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to initialize Postgres DB: %v", err)
	}

	redisClient, err := database.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis DB: %v", err)
	}

	repos := repository.NewRepository(postgresClient)
	redis := cache.NewCache(redisClient)

	tokenManager, err := auth.NewManager(cfg.JWT.SigningKey)
	if err != nil {
		log.Fatalf("Failed to initialize token manager: %v", err)
	}

	hasher := hash.NewSHA1Hasher("lolkek")

	serv := service.NewService(repos, logger, cfg.JWT, tokenManager, hasher, postgresClient, redis)

	handlers := http.NewHandler(serv, tokenManager)

	srv := server.NewServer(cfg.HTTP, handlers.Init())
	go func() {
		if err := srv.Run(); err != nil {
			logger.Info("shutting down server", slog.String("reason", err.Error()))
		}
	}()

	logger.Info("server started", slog.String("address", cfg.HTTP.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logger.Error("failed to stop server", slog.String("reason", err.Error()))
	}
}

// setupLogger initializes and returns a new logger instance configured
// with a text handler that outputs to the standard output.
// The logger is set to debug level and includes the source of the log.
func setupLogger() *slog.Logger {
	var logger *slog.Logger

	// Create a text handler for logging, outputting to standard output
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	logger = slog.New(handler)

	return logger
}
