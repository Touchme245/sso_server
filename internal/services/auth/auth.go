package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Touchme245/sso_server/internal/domain/models"
	"github.com/Touchme245/sso_server/internal/lib/jwt"
	"github.com/Touchme245/sso_server/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash string) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, uid int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

func NewAuthService(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, emal string, password string, appId int) (string, error) {
	const op = "authService.Login"
	log := slog.With(slog.String("op", op), slog.String("email", emal))
	log.Info("logining user")

	user, err := a.userProvider.User(ctx, emal)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", err)
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to find user", err)
		return "", fmt.Errorf("%s: %w", op, err)

	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil

}

func (a *AuthService) RegisterNewUser(ctx context.Context, emal string, password string) (int64, error) {
	const op = "aouthService.RegisterNewUser"

	log := slog.With(slog.String("op", op), slog.String("email", emal))
	log.Info("registering user")

	hasedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("Failed to generate password hash", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, emal, string(hasedPass))

	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return 0, ErrUserExists
		}
		log.Error("Failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil

}

func (a *AuthService) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "authService.Login"
	log := slog.With(slog.String("op", op))
	log.Info("is user admin check")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)

	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return false, ErrInvalidAppId
		}

		log.Error("Failed to check user is admin", err)
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked is user admin: ", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil

}
