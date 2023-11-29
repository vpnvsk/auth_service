package services

import (
	"auth_service/internal/models"
	"auth_service/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	AppId  int
}

func (s *Service) LoginUser(ctx context.Context, login, password string, appId int) (string, error) {
	const op = "services.LoginUser"

	user, err := s.repo.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.log.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		s.log.Info("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := s.repo.GetApp(ctx, appId)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := s.generateToken(user.Uuid, app)
	if err != nil {
		s.log.Error("failed to generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
func (s *Service) generateToken(userId uuid.UUID, app models.App) (string, error) {
	const op = "service.generateToken"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt: time.Now().Unix()},
		userId,
		app.ID,
	})
	return token.SignedString([]byte(app.Secret))

}

func (s *Service) RegisterUser(ctx context.Context, login, password string) (string, error) {
	const op = "service.RegisterUser"

	log := s.log.With(
		slog.String("op", op),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	id, err := s.repo.SaveUser(ctx, login, passHash)
	if err != nil {
		log.Error("failed to save user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id.String(), nil
}

func (s *Service) UserIsAdmin(ctx context.Context, userid string, appId int) (bool, error) {
	const op = "service.UserIsAdmin"
	log := s.log.With(
		slog.String("op", op),
	)
	userId, err := uuid.Parse(userid)
	if err != nil {
		log.Error("error while parsing userId", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}
	isAdmin, err := s.repo.IsAdmin(ctx, userId, appId)
	if err != nil {
		log.Error("failed to check if user is admin", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
