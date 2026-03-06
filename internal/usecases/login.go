package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/auth/jwt"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/auth/password"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/auth/token"
	"github.com/fedotovmax/mkk-luna-test/internal/config"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/errs"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/fedotovmax/mkk-luna-test/internal/ports"
	"github.com/fedotovmax/mkk-luna-test/internal/queries"
	"github.com/google/uuid"
)

type Login struct {
	log          *slog.Logger
	tokensCfg    *config.Tokens
	tokenManager ports.TokenManager
	users        queries.Users
	storage      ports.SessionStorage
}

func NewLogin(
	log *slog.Logger,
	users queries.Users,
	storage ports.SessionStorage,
	tokenManager ports.TokenManager,
	tokensCfg *config.Tokens,
) *Login {
	return &Login{
		log:          log,
		users:        users,
		storage:      storage,
		tokensCfg:    tokensCfg,
		tokenManager: tokenManager,
	}
}

func (u *Login) Execute(ctx context.Context, in *inputs.Login) (*domain.LoginResponse, error) {

	const op = "usecases.login"

	user, err := u.users.FindByEmail(ctx, in.Email)

	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrBadCredentials, err)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ok := password.ComparePasswords(in.Password, user.PasswordHash)

	if !ok {
		return nil, fmt.Errorf("%s: %w: %v", op, errs.ErrBadCredentials, err)
	}

	sid := uuid.New().String()

	nowUTC := time.Now().UTC()

	refreshToken, err := token.CreateToken()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	refreshExpTime := nowUTC.Add(u.tokensCfg.RefreshExpDuration)
	accessExpTime := nowUTC.Add(u.tokensCfg.AccessExpDuration)

	accessToken, err := u.tokenManager.Create(&jwt.CreateParams{
		Issuer:         u.tokensCfg.Issuer,
		Uid:            user.ID,
		Sid:            sid,
		TokenExpiresAt: accessExpTime,
		Now:            nowUTC,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = u.storage.Create(ctx, &inputs.CreateSession{
		ID:          sid,
		UserID:      user.ID,
		RefreshHash: refreshToken.Hashed,
		ExpiresAt:   refreshExpTime,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &domain.LoginResponse{
		AccessToken:    accessToken,
		RefreshToken:   refreshToken.Nohashed,
		AccessExpTime:  accessExpTime,
		RefreshExpTime: refreshExpTime,
	}, nil

}
