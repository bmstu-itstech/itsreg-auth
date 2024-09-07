package auth

import (
	"context"
	"errors"
	"fmt"
)

type UserNotFound struct {
	UserUUID string
}

func (e UserNotFound) Error() string {
	return fmt.Sprintf("user with UUID %s not found", e.UserUUID)
}

type UserEmailNotFound struct {
	Email string
}

func (e UserEmailNotFound) Error() string {
	return fmt.Sprintf("user with email %s not found", e.Email)
}

var ErrUserAlreadyExists = errors.New("user already exists")

type UsersRepository interface {
	Save(ctx context.Context, u *User) error
	User(ctx context.Context, uuid string) (*User, error)
	UserByEmail(ctx context.Context, email string) (*User, error)
	Update(
		ctx context.Context,
		uuid string,
		updateFn func(ctx context.Context, u *User) error,
	) error
	Delete(ctx context.Context, uuid string) error
}
