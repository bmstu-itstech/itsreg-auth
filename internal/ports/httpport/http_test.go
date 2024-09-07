package httpport_test

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-chi/chi/v5"
	"github.com/itsreg-auth/internal/common/jwtauth"
	"github.com/itsreg-auth/internal/common/server"
	"github.com/itsreg-auth/internal/common/tests"
	"github.com/itsreg-auth/internal/ports/httpport"
	"github.com/itsreg-auth/internal/service"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAuthHTTP(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping component test in short mode.")
	}

	t.Parallel()

	port := os.Getenv("PORT")
	addr := fmt.Sprintf("http://localhost:%s/api", port)
	client := httpport.MustNewHTTPAuthClient(addr)

	t.Run("should register user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		uuid := gofakeit.UUID()
		email := gofakeit.Email()
		password := fakePassword()

		res, err := client.RegisterUser(ctx, uuid, email, password)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		user, res, err := client.GetUser(ctx, uuid)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, res.StatusCode)
		require.Equal(t, email, user.Email)
		require.Equal(t, uuid, user.Uuid)
		require.Less(t, time.Now().Sub(user.CreatedAt), time.Second)
		require.Less(t, time.Now().Sub(user.UpdatedAt), time.Second)
	})

	t.Run("should login user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		uuid := gofakeit.UUID()
		email := gofakeit.Email()
		password := fakePassword()

		_, err := client.RegisterUser(ctx, uuid, email, password)
		require.NoError(t, err)

		tokens, res, err := client.LoginUser(ctx, email, password)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)

		parsed, err := jwtauth.ParseAccessToken(tokens.AccessToken)
		require.NoError(t, err)
		require.Equal(t, uuid, parsed.UserUUID)
	})

	t.Run("should return error if password mismatch", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		uuid := gofakeit.UUID()
		email := gofakeit.Email()
		password := fakePassword()

		_, err := client.RegisterUser(ctx, uuid, email, password)
		require.NoError(t, err)

		_, res, err := client.LoginUser(ctx, email, fakePassword())
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("should return error if user is invalid", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		uuid := gofakeit.UUID()
		email := ""
		password := fakePassword()

		res, err := client.RegisterUser(ctx, uuid, email, password)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("should return error if user not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		_, res, err := client.GetUser(ctx, gofakeit.UUID())
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

func fakePassword() string {
	return gofakeit.Password(true, true, true, true, false, 8)
}

func startService() bool {
	app := service.NewComponentTestApplication()

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	go server.RunHTTPServerOnAddr(addr, func(router chi.Router) http.Handler {
		return httpport.HandlerFromMux(httpport.NewHTTPServer(app), router)
	})

	ok := tests.WaitForPort(addr)
	if !ok {
		log.Println("Timed out waiting for auth HTTP to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
