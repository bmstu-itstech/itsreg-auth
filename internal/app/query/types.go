package query

import (
	"time"

	"github.com/itsreg-auth/internal/domain/auth"
)

type Empty struct{}

type User struct {
	UUID      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func mapUserFromDomain(u *auth.User) User {
	return User{
		UUID:      u.UUID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
