package infra_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/itsreg-auth/internal/domain/auth"
	"github.com/bmstu-itstech/itsreg-auth/internal/infra"
)

func TestPgUsersRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	url := os.Getenv("DATABASE_URI")
	db := sqlx.MustConnect("postgres", url)
	t.Cleanup(func() {
		err := db.Close()
		require.NoError(t, err)
	})

	repos := infra.NewPgUserRepository(db)
	testUsersRepository(t, repos)
}

func testUsersRepository(t *testing.T, r auth.UsersRepository) {
	t.Parallel()

	t.Run("should save user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		user := fakeUser()

		err := r.Save(ctx, user)
		require.NoError(t, err)

		saved, err := r.User(ctx, user.UUID)
		require.NoError(t, err)

		requireEqualUsers(t, user, saved)
	})

	t.Run("should return error if user uuid already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		user1 := fakeUser()
		err := r.Save(ctx, user1)
		require.NoError(t, err)

		user2 := fakeUser()
		user2.UUID = user1.UUID
		err = r.Save(ctx, user2)
		require.ErrorIs(t, err, auth.ErrUserAlreadyExists)
	})

	t.Run("should return error if user email already exists", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		user1 := fakeUser()
		err := r.Save(ctx, user1)
		require.NoError(t, err)

		user2 := fakeUser()
		user2.Email = user1.Email
		err = r.Save(ctx, user2)
		require.ErrorIs(t, err, auth.ErrUserAlreadyExists)
	})

	t.Run("should return error if user not found by uuid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		fakeUUID := gofakeit.UUID()
		_, err := r.User(ctx, fakeUUID)
		require.EqualError(t, err, fmt.Sprintf("user with UUID %s not found", fakeUUID))
	})

	t.Run("should return error if user not found by email", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		fakeEmail := gofakeit.Email()
		_, err := r.UserByEmail(ctx, fakeEmail)
		require.EqualError(t, err, fmt.Sprintf("user with email %s not found", fakeEmail))
	})

	t.Run("should update user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		user := fakeUser()
		err := r.Save(ctx, user)
		require.NoError(t, err)

		time.Sleep(time.Millisecond)

		err = r.Update(ctx, user.UUID, func(ctx context.Context, u *auth.User) error {
			u.UpdatedAt = time.Now()
			return nil
		})
		require.NoError(t, err)

		updated, err := r.User(ctx, user.UUID)
		require.NoError(t, err)

		require.Greater(t, updated.UpdatedAt.Sub(user.UpdatedAt), time.Millisecond)
	})

	t.Run("should return error on update if user not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		fakeUUID := gofakeit.UUID()
		err := r.Update(ctx, fakeUUID, func(ctx context.Context, u *auth.User) error {
			return nil
		})
		require.EqualError(t, err, fmt.Sprintf("user with UUID %s not found", fakeUUID))
	})

	t.Run("should delete user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		user := fakeUser()
		err := r.Save(ctx, user)

		err = r.Delete(ctx, user.UUID)
		require.NoError(t, err)

		_, err = r.User(ctx, user.UUID)
		require.EqualError(t, err, fmt.Sprintf("user with UUID %s not found", user.UUID))
	})

	t.Run("should return error if user for delete not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		fakeUUID := gofakeit.UUID()
		err := r.Delete(ctx, fakeUUID)
		require.EqualError(t, err, fmt.Sprintf("user with UUID %s not found", fakeUUID))
	})
}

const minPasswordLength = 8

func fakePassword() string {
	return gofakeit.Password(true, true, true, true, false, minPasswordLength)
}

func fakeUser() *auth.User {
	return auth.MustNewUser(
		gofakeit.UUID(),
		gofakeit.Email(),
		fakePassword(),
	)
}

func requireEqualUsers(t *testing.T, expected *auth.User, actual *auth.User) {
	require.True(
		t, equalUsers(*expected, *actual),
		fmt.Sprintf(
			"expected and actual users are different\nexpected: %+v\nactual: %+v",
			expected, actual,
		),
	)
}

func equalUsers(a auth.User, b auth.User) bool {
	return a.UUID == b.UUID &&
		a.Email == b.Email &&
		bytes.Compare(a.Passhash, b.Passhash) == 0 &&
		a.CreatedAt.Sub(b.CreatedAt) < time.Microsecond &&
		a.UpdatedAt.Sub(b.UpdatedAt) < time.Microsecond
}
