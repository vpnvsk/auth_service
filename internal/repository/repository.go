package repository

import (
	"auth_service/internal/models"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type Repository interface {
	Saver
	Getter
}
type Repo struct {
	log    *slog.Logger
	saver  Saver
	getter Getter
}

func New(log *slog.Logger,
	saver Saver,
	getter Getter) *Repo {
	return &Repo{log: log, saver: saver, getter: getter}
}

type Saver interface {
	SaveUser(ctx context.Context, user models.User) (uuid.UUID, error)
}

type Getter interface {
	GetUser(ctx context.Context, login, password string, addId int) (models.User, error)
	IsAdmin(ctx context.Context, uuid string) (bool, error)
}
