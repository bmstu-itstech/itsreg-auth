package auth

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/itsreg-auth/internal/common/commonerrs"
)

type User struct {
	UUID string

	Email    string
	Passhash []byte

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(
	uuid string,
	email string,
	password string,
) (*User, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty uuid")
	}

	if email == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty email")
	}

	if password == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty password")
	}

	passhash, err := createPasshash(password)
	if err != nil {
		return nil, err
	}

	return &User{
		UUID:      uuid,
		Email:     email,
		Passhash:  passhash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func MustNewUser(
	uuid string,
	email string,
	password string,
) *User {
	user, err := NewUser(uuid, email, password)
	if err != nil {
		panic(err)
	}
	return user
}

func NewUserFromDB(
	uuid string,
	email string,
	passhash []byte,
	createdAt time.Time,
	updatedAt time.Time,
) (*User, error) {
	if uuid == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty uuid")
	}

	if email == "" {
		return nil, commonerrs.NewInvalidInputError("expected not empty email")
	}

	if len(passhash) == 0 {
		return nil, commonerrs.NewInvalidInputError("expected not empty password")
	}

	if createdAt.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not empty createdAt")
	}

	if updatedAt.IsZero() {
		return nil, commonerrs.NewInvalidInputError("expected not empty updatedAt")
	}

	return &User{
		UUID:      uuid,
		Email:     email,
		Passhash:  passhash,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (u *User) PasswordMatch(password string) error {
	if err := bcrypt.CompareHashAndPassword(u.Passhash, []byte(password)); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

func createPasshash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
