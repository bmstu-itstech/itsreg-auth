package value

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Email string

func ValidateEmail(email Email) error {
	return validation.Validate(
		email,
		validation.Required,
		is.Email,
	)
}
