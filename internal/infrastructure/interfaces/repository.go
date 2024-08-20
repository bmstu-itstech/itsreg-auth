package interfaces

import (
	"context"
	"errors"

	"github.com/itsreg-auth/internal/domain/entity"
	"github.com/itsreg-auth/internal/domain/value"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) (value.UserId, error)
	Find(ctx context.Context, user *entity.User) (*entity.User, error)
}
