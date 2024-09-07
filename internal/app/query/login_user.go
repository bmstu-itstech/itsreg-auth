package query

import (
	"context"
	"errors"
	"log/slog"

	"github.com/itsreg-auth/internal/common/decorator"
	"github.com/itsreg-auth/internal/domain/auth"
)

type LoginUser struct {
	Email    string
	Password string
}

type LoginUserHandler decorator.QueryHandler[LoginUser, User]

type loginUserHandler struct {
	users auth.UsersRepository
}

func NewLoginUserHandler(
	users auth.UsersRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) LoginUserHandler {
	if users == nil {
		panic("users repository is nil")
	}

	return decorator.ApplyQueryDecorators[LoginUser, User](
		loginUserHandler{users: users},
		logger,
		metricsClient,
	)
}

func (h loginUserHandler) Handle(ctx context.Context, query LoginUser) (User, error) {
	user, err := h.users.UserByEmail(ctx, query.Email)
	if errors.As(err, &auth.UserEmailNotFound{}) {
		return User{}, auth.ErrInvalidCredentials
	} else if err != nil {
		return User{}, err
	}

	if err = user.PasswordMatch(query.Password); err != nil {
		return User{}, err
	}

	return mapUserFromDomain(user), nil
}
