package render

import (
	"bytes"
	"log"
	"net/http"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
)

var App *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	App = a
}

// AddDefaultData adds default data to template data
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = App.Session.PopString(r.Context(), "flash")
	td.Error = App.Session.PopString(r.Context(), "error")
	td.Warning = App.Session.PopString(r.Context(), "warning")
	return td
}

// ---------- PUBLIC PAGES ----------
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	ts, ok := App.TemplateCache[tmpl]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", AddDefaultData(td, r))
	if err != nil {
		log.Println(err)
	}

	_, _ = buf.WriteTo(w)
}

// ---------- ADMIN PAGES ----------
func AdminTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	ts, ok := App.TemplateCache[tmpl]
	if !ok {
		http.Error(w, "admin template not found", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "admin.layout", AddDefaultData(td, r))
	if err != nil {
		log.Println(err)
	}

	_, _ = buf.WriteTo(w)
}
