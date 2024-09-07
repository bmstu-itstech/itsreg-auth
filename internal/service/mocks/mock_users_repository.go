package mocks

import (
	"context"
	"sync"

	"github.com/itsreg-auth/internal/domain/auth"
)

type mockUserRepository struct {
	sync.RWMutex
	m map[string]auth.User
}

func NewMockUserRepository() auth.UsersRepository {
	return &mockUserRepository{
		m: make(map[string]auth.User),
	}
}

func (r *mockUserRepository) Save(ctx context.Context, u *auth.User) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[u.UUID]; ok {
		return auth.ErrUserAlreadyExists
	}

	for _, user := range r.m {
		if user.Email == u.Email {
			return auth.ErrUserAlreadyExists
		}
	}

	r.m[u.UUID] = *u

	return nil
}

func (r *mockUserRepository) User(ctx context.Context, uuid string) (*auth.User, error) {
	r.RLock()
	defer r.RUnlock()

	u, ok := r.m[uuid]
	if !ok {
		return nil, auth.UserNotFound{UserUUID: uuid}
	}

	return &u, nil
}

func (r *mockUserRepository) UserByEmail(ctx context.Context, email string) (*auth.User, error) {
	r.RLock()
	defer r.RUnlock()

	for _, u := range r.m {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, auth.UserEmailNotFound{Email: email}
}

func (r *mockUserRepository) Update(
	ctx context.Context,
	uuid string,
	updateFn func(ctx context.Context, u *auth.User) error,
) error {
	r.Lock()
	defer r.Unlock()

	user, ok := r.m[uuid]
	if !ok {
		return auth.UserNotFound{UserUUID: uuid}
	}

	err := updateFn(ctx, &user)
	if err != nil {
		return err
	}

	r.m[uuid] = user

	return nil
}

func (r *mockUserRepository) Delete(ctx context.Context, uuid string) error {
	r.Lock()
	defer r.Unlock()

	_, ok := r.m[uuid]
	if !ok {
		return auth.UserNotFound{UserUUID: uuid}
	}

	delete(r.m, uuid)

	return nil
}
