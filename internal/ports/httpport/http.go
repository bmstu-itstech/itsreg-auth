package httpport

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/render"

	"github.com/bmstu-itstech/itsreg-auth/internal/app"
	"github.com/bmstu-itstech/itsreg-auth/internal/app/command"
	"github.com/bmstu-itstech/itsreg-auth/internal/app/query"
	"github.com/bmstu-itstech/itsreg-auth/internal/common/commonerrs"
	"github.com/bmstu-itstech/itsreg-auth/internal/common/jwtauth"
	"github.com/bmstu-itstech/itsreg-auth/internal/domain/auth"
)

const (
	accessTTL = time.Hour * 24 * 7
)

type Server struct {
	app *app.Application
}

func NewHTTPServer(app *app.Application) *Server {
	return &Server{app: app}
}

func (s Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var postRegister PostRegister
	if err := render.Decode(r, &postRegister); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	err := s.app.Commands.RegisterUser.Handle(r.Context(), command.RegisterUser{
		UUID:     postRegister.Uuid,
		Email:    postRegister.Email,
		Password: postRegister.Password,
	})
	if errors.As(err, &commonerrs.InvalidInputError{}) {
		httpError(w, r, err, http.StatusBadRequest)
		return
	} else if errors.Is(err, auth.ErrUserAlreadyExists) {
		httpError(w, r, err, http.StatusBadRequest)
		return
	} else if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-location", fmt.Sprintf("/users/%s", postRegister.Uuid))
	w.WriteHeader(http.StatusCreated)
}

func (s Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	var postLogin PostLogin
	if err := render.Decode(r, &postLogin); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}

	user, err := s.app.Queries.LoginUser.Handle(r.Context(), query.LoginUser{
		Email:    postLogin.Email,
		Password: postLogin.Password,
	})
	if errors.Is(err, auth.ErrInvalidCredentials) {
		httpError(w, r, err, http.StatusUnauthorized)
		return
	} else if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	at, err := jwtauth.NewAccessToken(user.UUID, accessTTL)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	res := Authenticated{
		AccessToken: at,
	}

	render.JSON(w, r, res)
}

func (s Server) GetUser(w http.ResponseWriter, r *http.Request, uuid string) {
	user, err := s.app.Queries.GetUser.Handle(r.Context(), query.GetUser{
		UserUUID: uuid,
	})
	if errors.As(err, &auth.UserNotFound{}) {
		httpError(w, r, err, http.StatusNotFound)
		return
	} else if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, mapUserToAPI(user))
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	w.WriteHeader(code)
	render.JSON(w, r, Error{Message: err.Error()})
}

func mapUserToAPI(user query.User) User {
	return User{
		CreatedAt: user.CreatedAt,
		Email:     user.Email,
		UpdatedAt: user.UpdatedAt,
		Uuid:      user.UUID,
	}
}
