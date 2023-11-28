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
	UserIsAdmin(ctx context.Context, uuid string) (bool, error)
}

type Service struct {
	log      *slog.Logger
	repo     repository.Repository
	tokenTTL time.Duration
}

func (*Service) LoginUser(ctx context.Context, login, password string, appId int) (string, error) {
	return "", nil
}

func (*Service) RegisterUser(ctx context.Context, login, password string) (string, error) {
	return "", nil
}

func (*Service) UserIsAdmin(ctx context.Context, uuid string) (bool, error) {
	return true, nil
}
