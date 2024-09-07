package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/itsreg-auth/internal/common/server"
	"github.com/itsreg-auth/internal/ports/httpport"
	"github.com/itsreg-auth/internal/service"
)

func main() {
	app, cleanup := service.NewApplication()
	defer cleanup()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return httpport.HandlerFromMux(httpport.NewHTTPServer(app), router)
	})
}
