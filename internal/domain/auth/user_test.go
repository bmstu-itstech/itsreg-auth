package auth_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-auth/internal/domain/auth"
)

func TestUser_MatchPassword(t *testing.T) {
	password := "qwerty"
	user := auth.MustNewUser("1234", "test@test.com", password)
	require.NoError(t, user.PasswordMatch(password))
	require.ErrorIs(t, user.PasswordMatch("another"), auth.ErrInvalidCredentials)
}
