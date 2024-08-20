package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/itsreg-auth/internal/domain/entity"
	"github.com/itsreg-auth/internal/domain/value"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) (value.UserId, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(value.UserId), args.Error(1)
}

func (m *MockUserRepository) Find(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*entity.User), args.Error(1)
}
