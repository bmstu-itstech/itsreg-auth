package app

import (
	"github.com/itsreg-auth/internal/app/command"
	"github.com/itsreg-auth/internal/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	RegisterUser command.RegisterUserHandler
}

type Queries struct {
	LoginUser query.LoginUserHandler
	GetUser   query.GetUserHandler
}
