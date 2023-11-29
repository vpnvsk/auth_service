package services

import (
	"auth_service/internal/repository"
	"context"
	"log/slog"
	"time"
)

type Services interface {
	LoginUser(ctx context.Context, login, password string, appId int) (string, error)
	RegisterUser(ctx context.Context, login, password string) (string, error)
	UserIsAdmin(ctx context.Context, uuid string, appId int) (bool, error)
}

type Service struct {
	log      *slog.Logger
	repo     repository.Repository
	tokenTTL time.Duration
}

func New(logger *slog.Logger, repo repository.Repository, tokenTTL time.Duration) *Service {
	return &Service{
		log:      logger,
		repo:     repo,
		tokenTTL: tokenTTL,
	}
}
