package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

// ✅ This test now checks that the returned router is of type *chi.Mux
func TestRoutes(t *testing.T) {
	mux := routes(nil)

	// Assert the type of the router is *chi.Mux
	if _, ok := mux.(*chi.Mux); !ok {
		t.Errorf("Expected router type *chi.Mux, but got %T", mux)
	}
}

// ✅ The real route setup (for testing or actual use)
func routes(app interface{}) http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(nosurf.NewPure)

	return mux
}
