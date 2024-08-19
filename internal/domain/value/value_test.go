package value

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUserId(t *testing.T) {
	tests := []struct {
		userId  UserId
		wantErr bool
	}{
		{0, false},
		{1, false},
		{^UserId(0), false}, // maximum UserId
		{UserId(^uint64(0)), false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("userId_%d", tt.userId), func(t *testing.T) {
			err := ValidateUserId(tt.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   Email
		wantErr bool
	}{
		{"", true},
		{"invalid-email", true},
		{"user@example.com", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password Password
		wantErr  bool
	}{
		{"", true},
		{"short", true},
		{"noDigits", true},
		{"123", true},
		{"Password1", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.password), func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPasswordHash(t *testing.T) {
	tests := []struct {
		password Password
		wantErr  bool
	}{
		{"Password1", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.password), func(t *testing.T) {
			hash, err := GetPasswordHash(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	password := Password("Password1")
	hash, err := GetPasswordHash(password)
	assert.NoError(t, err)

	tests := []struct {
		hash     PasswordHash
		password Password
		want     bool
		wantErr  bool
	}{
		{hash, password, true, false},
		{hash, "wrongPassword", false, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.password), func(t *testing.T) {
			match, err := CompareHashAndPassword(tt.hash, tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, match)
			}
		})
	}
}
