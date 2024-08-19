package entity

import (
	"errors"

	"github.com/itsreg-auth/internal/domain/value"
)

type User struct {
	Id           value.UserId
	Email        value.Email
	PasswordHash value.PasswordHash
}

var ErrDefaultUserId = errors.New("default zero user id")

// NewUserFromDB Constructor for creating a user from the database
func NewUserFromDB(id value.UserId, email value.Email, passwordHash value.PasswordHash) (User, error) {
	if err := value.ValidateUserId(id); err != nil {
		return User{}, err
	}
	return User{
		Id:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

// NewUserRegistration Constructor for user registration
func NewUserRegistration(email value.Email, password value.Password) (User, error) {
	if err := value.ValidateEmail(email); err != nil {
		return User{}, err
	}
	if err := value.ValidatePassword(password); err != nil {
		return User{}, err
	}
	passwordHash, err := value.GetPasswordHash(password)
	if err != nil {
		return User{}, err
	}
	return User{
		Id:           value.UnknownUserId,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}

func SetUserId(user User, id value.UserId) (User, error) {
	if user.Id != value.UnknownUserId {
		return User{}, ErrDefaultUserId
	}
	user.Id = id
	return user, nil
}
