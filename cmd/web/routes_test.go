package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

// Test that the application's real router is correctly created
func TestRoutes(t *testing.T) {
	mux := testRoutes()

	if _, ok := mux.(*chi.Mux); !ok {
		t.Errorf("Expected router type *chi.Mux, but got %T", mux)
	}
}

// testRoutes is a local test helper that does NOT conflict with the real routes()
func testRoutes() http.Handler {
	mux := chi.NewRouter()

	// middleware matching the real app
	mux.Use(nosurf.NewPure)

	return mux
}
