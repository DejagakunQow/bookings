package main

import (
	"net/http"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// routes sets up the application routes and middleware
func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// --- Middleware ---
	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad(app)) // REQUIRED for SCS sessions

	// --- Public routes ---
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)

	// --- Static files ---
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
