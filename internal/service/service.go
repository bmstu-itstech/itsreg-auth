package service

import (
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/itsreg-auth/internal/app"
	"github.com/itsreg-auth/internal/app/command"
	"github.com/itsreg-auth/internal/app/query"
	"github.com/itsreg-auth/internal/common/decorator"
	"github.com/itsreg-auth/internal/common/logs"
	"github.com/itsreg-auth/internal/common/metrics"
	"github.com/itsreg-auth/internal/domain/auth"
	"github.com/itsreg-auth/internal/infra"
	"github.com/itsreg-auth/internal/service/mocks"
)

type Cleanup func()

func NewApplication() (*app.Application, Cleanup) {
	logger := logs.DefaultLogger()
	metricsClient := metrics.NoOp{}

	url := os.Getenv("DATABASE_URI")
	db := sqlx.MustConnect("postgres", url)

	users := infra.NewPgUserRepository(db)

	return newApplication(logger, metricsClient, users), func() {
		_ = db.Close()
	}
}

func NewComponentTestApplication() *app.Application {
	logger := logs.DefaultLogger()
	metricsClient := metrics.NoOp{}

	users := mocks.NewMockUserRepository()

	return newApplication(logger, metricsClient, users)
}

func newApplication(
	logger *slog.Logger,
	metricsClients decorator.MetricsClient,
	users auth.UsersRepository,
) *app.Application {
	return &app.Application{
		Commands: app.Commands{
			RegisterUser: command.NewRegisterUserHandler(users, logger, metricsClients),
		},
		Queries: app.Queries{
			GetUser:   query.NewGetUserHandler(users, logger, metricsClients),
			LoginUser: query.NewLoginUserHandler(users, logger, metricsClients),
		},
	}
}
