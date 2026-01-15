package main

import (
	"net/http"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/go-chi/chi"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// TEMP: no middleware, no sessions, no CSRF
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return mux
}
