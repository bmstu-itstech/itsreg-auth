package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/itsreg-auth/internal/common/decorator"
	"github.com/bmstu-itstech/itsreg-auth/internal/domain/auth"
)

type GetUser struct {
	UserUUID string
}

type GetUserHandler decorator.QueryHandler[GetUser, User]

type getUserHandler struct {
	users auth.UsersRepository
}

func NewGetUserHandler(
	users auth.UsersRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) GetUserHandler {
	if users == nil {
		panic("users repository is nil")
	}

	return decorator.ApplyQueryDecorators[GetUser, User](
		getUserHandler{users: users},
		logger,
		metricsClient,
	)
}

func (h getUserHandler) Handle(ctx context.Context, query GetUser) (User, error) {
	user, err := h.users.User(ctx, query.UserUUID)
	if err != nil {
		return User{}, err
	}

	return mapUserFromDomain(user), nil
}
