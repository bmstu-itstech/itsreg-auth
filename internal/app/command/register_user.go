package command

import (
	"context"
	"log/slog"

	"github.com/itsreg-auth/internal/common/decorator"
	"github.com/itsreg-auth/internal/domain/auth"
)

type RegisterUser struct {
	UUID     string
	Email    string
	Password string
}

type RegisterUserHandler decorator.CommandHandler[RegisterUser]

type registerUserHandler struct {
	users auth.UsersRepository
}

func NewRegisterUserHandler(
	users auth.UsersRepository,

	logger *slog.Logger,
	metricsClient decorator.MetricsClient,
) RegisterUserHandler {
	if users == nil {
		panic("users repository is nil")
	}

	return decorator.ApplyCommandDecorators[RegisterUser](
		&registerUserHandler{users: users},
		logger,
		metricsClient,
	)
}

func (h registerUserHandler) Handle(ctx context.Context, cmd RegisterUser) error {
	user, err := auth.NewUser(cmd.UUID, cmd.Email, cmd.Password)
	if err != nil {
		return err
	}

	return h.users.Save(ctx, user)
}
