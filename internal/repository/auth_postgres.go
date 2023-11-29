package repository

import (
	"auth_service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type DBSaver struct {
	log *slog.Logger
	db  *sqlx.DB
}
type DBReader struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewDBSaver(logger *slog.Logger, db *sqlx.DB) *DBSaver {
	return &DBSaver{
		log: logger,
		db:  db,
	}
}
func NewDBReader(logger *slog.Logger, db *sqlx.DB) *DBReader {
	return &DBReader{
		log: logger,
		db:  db,
	}
}

func (r *DBSaver) SaveUser(ctx context.Context, login string, passwordHash []byte) (uuid.UUID, error) {
	const op = "repo.auth_postgres.SaveUser"
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (login, password_hash) VALUES ($1, $2) RETURNING id", userTable)
	row := r.db.QueryRow(query, login, passwordHash)
	if err := row.Scan(&id); err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (r *DBReader) GetUser(ctx context.Context, login string) (models.User, error) {
	const op = "repo.auth_postgres.GetUser"
	var user models.User
	query := fmt.Sprintf("SELECT id, login, password_hash FROM %s WHERE login=$1", userTable)
	if err := r.db.Get(&user, query, login); err != nil {
		//TODO: specify error
		return user, fmt.Errorf("%s: %w", op, err)
		//return user, ErrUserNotFound
	}

	return user, nil
}

func (r *DBReader) IsAdmin(ctx context.Context, uuid uuid.UUID, appId int) (bool, error) {
	const op = "repo.auth_postgres.IsAdmin"
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id=$1 AND app_id=$2", adminTable)
	row := r.db.QueryRow(query, uuid, appId)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return count > 0, nil
}
func (r *DBReader) GetApp(ctx context.Context, id int) (models.App, error) {
	const op = "repo.auth_postgres.GetApp"
	var app models.App
	query := fmt.Sprintf("SELECT id, name, secret FROM %s WHERE id=$1", appTable)
	if err := r.db.Get(&app, query, id); err != nil {
		//TODO: specify error
		return app, err
	}
	return app, nil
}
