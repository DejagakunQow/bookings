package render

import (
	"bytes"
	"log"
	"net/http"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
)

var app *config.AppConfig

// NewRenderer sets the application config for the render package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds common data to all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	if td == nil {
		td = &models.TemplateData{}
	}

	// SAFETY: if session is not initialized, do NOT panic
	if app == nil || app.Session == nil {
		return td
	}

	// Flash messages
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	// Authentication flag
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = true
	}

	return td
}

// ---------- PUBLIC PAGES ----------
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	if app == nil {
		http.Error(w, "Application not initialized", http.StatusInternalServerError)
		return
	}

	ts, ok := app.TemplateCache[tmpl]
	if !ok {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", AddDefaultData(td, r))
	if err != nil {
		log.Println(err)
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}

	_, _ = buf.WriteTo(w)
}

// ---------- ADMIN PAGES ----------
func AdminTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	if app == nil {
		http.Error(w, "Application not initialized", http.StatusInternalServerError)
		return
	}

	ts, ok := app.TemplateCache[tmpl]
	if !ok {
		http.Error(w, "Admin template not found", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "admin.layout", AddDefaultData(td, r))
	if err != nil {
		log.Println(err)
		http.Error(w, "Admin template execution error", http.StatusInternalServerError)
		return
	}

	_, _ = buf.WriteTo(w)
}
