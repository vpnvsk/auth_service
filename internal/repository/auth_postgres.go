package repository

import (
	"auth_service/internal/models"
	"context"
	"github.com/google/uuid"
	"log/slog"
)

type Save struct {
	log *slog.Logger
}

func (r *Save) SaveUser(ctx context.Context, user models.User) (uuid.UUID, error) {
	return uuid.Nil, nil
}
