package render

import (
	"bytes"
	"log"
	"net/http"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
)

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds default data to template data
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	if td == nil {
		td = &models.TemplateData{}
	}
	return td
}

// ---------- PUBLIC PAGES ----------
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	ts, ok := app.TemplateCache[tmpl]
	if !ok {
		log.Fatal("template not found in cache: ", tmpl)
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
	ts, ok := app.TemplateCache[tmpl]
	if !ok {
		log.Fatal("admin template not found in cache: ", tmpl)
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "admin.layout", AddDefaultData(td, r))
	if err != nil {
		log.Println(err)
	}

	_, _ = buf.WriteTo(w)
}
