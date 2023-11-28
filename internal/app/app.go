package app

import (
	"auth_service/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPSServer *grpc.App
}

func New(log *slog.Logger, port int, tokenTTL time.Duration) *App {
	// TODO: init storage

	// TODO: init service

	grpcApp := grpc.New(log, port)
	return &App{
		GRPSServer: grpcApp,
	}
}
