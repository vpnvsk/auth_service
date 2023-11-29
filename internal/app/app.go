package app

import (
	"auth_service/internal/app/grpc"
	"auth_service/internal/repository"
	"auth_service/internal/services"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"time"
)

type App struct {
	GRPSServer *grpc.App
}

func New(log *slog.Logger, port int, tokenTTL time.Duration) *App {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	db, err := repository.NewPostgresDb(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		panic(err)
	}
	repo := repository.New(log, db)
	service := services.New(log, repo, tokenTTL)
	grpcApp := grpc.New(log, port, service)
	return &App{
		GRPSServer: grpcApp,
	}
}
