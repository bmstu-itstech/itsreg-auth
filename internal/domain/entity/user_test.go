package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itsreg-auth/internal/domain/value"
)

func TestNewUserFromDB(t *testing.T) {
	tests := []struct {
		id           value.UserId
		email        value.Email
		passwordHash value.PasswordHash
		wantErr      bool
	}{
		{1, "user@example.com", []byte("$2a$10$7bH1ZZP8fP2ROXxkH1H43eMzXO2R5R/0CR6Pqf/7Q6J5q4P1/dxZC"), false},
		{0, "user@example.com", []byte("$2a$10$7bH1ZZP8fP2ROXxkH1H43eMzXO2R5R/0CR6Pqf/7Q6J5q4P1/dxZC"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			user, err := NewUserFromDB(tt.id, tt.email, tt.passwordHash)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
			}
		})
	}
}

func TestNewUserRegistration(t *testing.T) {
	tests := []struct {
		email    value.Email
		password value.Password
		wantErr  bool
	}{
		{"user@example.com", "Password1", false},
		{"invalid-email", "Password1", true},
		{"user@example.com", "short", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			user, err := NewUserRegistration(tt.email, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.NotEqual(t, user.PasswordHash, tt.password)
			}
		})
	}
}
