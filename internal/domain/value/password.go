package value

import (
	"errors"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	Password     string
	PasswordHash []byte
)

func ValidatePassword(password Password) error {
	return validation.Validate(
		password,
		validation.Required,
		validation.Length(8, 0),
		validation.Match(regexp.MustCompile(`[a-zA-Z]`)).Error("must contain at least one letter"),
		validation.Match(regexp.MustCompile(`[0-9]`)).Error("must contain at least one digit"),
	)
}

func GetPasswordHash(password Password) (PasswordHash, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to get password hash %w", err)
	}
	return passwordHash, nil
}

func CompareHashAndPassword(hash PasswordHash, password Password) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return true, nil
}
