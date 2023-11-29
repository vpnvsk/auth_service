package repository

import (
	"auth_service/internal/models"
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Repository interface {
	Saver
	Getter
}
type Repo struct {
	log *slog.Logger
	Saver
	Getter
}

func New(log *slog.Logger, db *sqlx.DB) *Repo {
	return &Repo{log: log, Saver: NewDBSaver(log, db), Getter: NewDBReader(log, db)}
}

type Saver interface {
	SaveUser(ctx context.Context, login string, passwordHash []byte) (uuid.UUID, error)
}

type Getter interface {
	GetUser(ctx context.Context, login string) (models.User, error)
	GetApp(ctx context.Context, id int) (models.App, error)
	IsAdmin(ctx context.Context, uuid uuid.UUID, appId int) (bool, error)
}
