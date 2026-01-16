package render

import (
	"html/template"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
)

var pathToTemplates string

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	Template(&ww, r, "home.page.tmpl", &models.TemplateData{})

	Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	if td == nil {
		td = &models.TemplateData{}
	}

	td.Flash = session.GetString(r.Context(), "flash")
	td.Error = session.GetString(r.Context(), "error")
	td.Warning = session.GetString(r.Context(), "warning")

	return td
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// All page templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return cache, err
	}

	// Public layouts (base)
	layouts, err := filepath.Glob("./templates/*.layout.tmpl")
	if err != nil {
		return cache, err
	}

	// Admin layouts
	adminLayouts, err := filepath.Glob("./templates/admin/*.layout.tmpl")
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		// Parse public layout
		ts, err = ts.ParseFiles(layouts...)
		if err != nil {
			return cache, err
		}

		// Parse admin layout
		ts, err = ts.ParseFiles(adminLayouts...)
		if err != nil {
			return cache, err
		}

		cache[name] = ts
	}

	return cache, nil
}
