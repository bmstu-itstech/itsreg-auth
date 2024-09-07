package httpport

import (
	"context"
	"github.com/go-chi/render"
	"github.com/itsreg-auth/internal/common/client/auth"
	"net/http"
)

type HTTPAuthClient struct {
	client *auth.Client
}

func NewHTTPAuthClient(addr string) (*HTTPAuthClient, error) {
	client, err := auth.NewClient(addr)
	if err != nil {
		return nil, err
	}
	return &HTTPAuthClient{
		client: client,
	}, nil
}

func MustNewHTTPAuthClient(addr string) *HTTPAuthClient {
	c, err := NewHTTPAuthClient(addr)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *HTTPAuthClient) RegisterUser(ctx context.Context, uuid string, email string, password string) (*http.Response, error) {
	return c.client.RegisterUser(ctx, auth.RegisterUserJSONRequestBody{
		Uuid:     uuid,
		Email:    email,
		Password: password,
	})
}

func (c *HTTPAuthClient) LoginUser(ctx context.Context, email string, password string) (auth.Authenticated, *http.Response, error) {
	res, err := c.client.LoginUser(ctx, auth.LoginUserJSONRequestBody{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return auth.Authenticated{}, res, err
	}

	var token auth.Authenticated
	if err = render.DecodeJSON(res.Body, &token); err != nil {
		return auth.Authenticated{}, res, err
	}

	return token, res, nil
}

func (c *HTTPAuthClient) GetUser(ctx context.Context, uuid string) (auth.User, *http.Response, error) {
	res, err := c.client.GetUser(ctx, uuid)
	if err != nil {
		return auth.User{}, res, err
	}

	var user auth.User
	if err = render.DecodeJSON(res.Body, &user); err != nil {
		return auth.User{}, res, err
	}

	return user, res, nil
}
