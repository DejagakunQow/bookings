package main

import (
	"net/http"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/handlers"
	"github.com/go-chi/chi"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// TEMP: no middleware (sessions / CSRF disabled for now)
	mux.Get("/", handlers.Repo.Home)

	// Static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
