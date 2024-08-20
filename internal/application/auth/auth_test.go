package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/itsreg-auth/internal/domain/entity"
	"github.com/itsreg-auth/internal/domain/value"
	"github.com/itsreg-auth/internal/infrastructure/interfaces"
	"github.com/itsreg-auth/internal/infrastructure/mocks"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo, "test_secret")

	t.Run("successful registration", func(t *testing.T) {
		user, err := entity.NewUserRegistration("test@example.com", "password123")
		assert.NoError(t, err, "failed create user")

		mockRepo.On("Find", mock.Anything, mock.Anything).Return(&user, interfaces.ErrUserNotFound).Once()
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(value.UserId(1), nil).Once()

		userId, err := service.Register(ctx, "test@example.com", "StrongPass123!")
		assert.NoError(t, err)
		assert.Equal(t, value.UserId(1), userId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("email already used", func(t *testing.T) {
		user, err := entity.NewUserRegistration("test@example.com", "password123")
		assert.NoError(t, err, "failed create user")

		mockRepo.On("Find", mock.Anything, mock.Anything).Return(&user, nil).Once()

		userId, err := service.Register(ctx, "test@example.com", "StrongPass123!")
		assert.ErrorIs(t, err, ErrEmailAlreadyUsed)
		assert.Equal(t, value.UnknownUserId, userId)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		userId, err := service.Register(ctx, "invalid-email", "short")
		assert.ErrorIs(t, err, ErrValidationUser)
		assert.Equal(t, value.UnknownUserId, userId)
	})

	t.Run("internal error during save", func(t *testing.T) {
		var user entity.User
		mockRepo.On("Find", mock.Anything, mock.Anything).Return(&user, interfaces.ErrUserNotFound).Once()
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(value.UnknownUserId, errors.New("db error")).Once()

		userId, err := service.Register(ctx, "test@example.com", "StrongPass123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AuthService.Register")
		assert.Equal(t, value.UnknownUserId, userId)
		mockRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mocks.MockUserRepository)
	service := NewAuthService(mockRepo, "test_secret")

	t.Run("successful login", func(t *testing.T) {
		user, err := entity.NewUserRegistration("test@example.com", "StrongPass123!")
		assert.NoError(t, err, "failed create user")

		mockRepo.On("Find", mock.Anything, mock.Anything).Return(&user, nil).Once()

		token, err := service.Login(ctx, "test@example.com", "StrongPass123!")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		var user entity.User
		mockRepo.On("Find", mock.Anything, mock.Anything).Return(&user, interfaces.ErrUserNotFound).Once()

		token, err := service.Login(ctx, "test@example.com", "StrongPass123!")
		assert.ErrorIs(t, err, ErrUserNotFound)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		user, err := entity.NewUserRegistration("test@example.com", "OtherPassword123!")
		assert.NoError(t, err, "failed create user")

		mockRepo.On("Find", mock.Anything, mock.Anything).Return(&user, nil).Once()

		token, err := service.Login(ctx, "test@example.com", "WrongPassword123!")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("internal error during find", func(t *testing.T) {
		var user entity.User
		mockRepo.On("Find", mock.Anything, mock.Anything).Return(&user, errors.New("db error")).Once()

		token, err := service.Login(ctx, "test@example.com", "StrongPass123!")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "AuthService.Login")
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
	})
}
