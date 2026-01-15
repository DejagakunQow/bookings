package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormartDate,
	"isoDate":    isoDate,
}

func isoDate(t time.Time) string {
	return t.Format("2006-01-02")
}

var app *config.AppConfig
var pathToTemplates = "./templates"

// HumanDate formats a time.Time into a human-readable date string
func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("02 Jan 2006")
}

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}
func FormartDate(t time.Time, f string) string {
	return t.Format(f)
}

// AddDefaultDataToTemplate adds default data for all templates
func AddDefaultDataToTemplate(td *models.TemplateData, r *http.Request) *models.TemplateData {
	if td == nil {
		td = &models.TemplateData{}
	}

	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = true

	return td
}

// Template renders templates
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("cannot get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultDataToTemplate(td, r)

	layout := "base"

	if strings.HasPrefix(tmpl, "admin-") {
		layout = "admin.layout"
	}

	err := t.ExecuteTemplate(buf, layout, td)
	if err != nil {
		log.Println("Template execution error:", err)
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache builds the template cache
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(pathToTemplates + "/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	layouts, err := filepath.Glob(pathToTemplates + "/*.layout.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := append([]string{page}, layouts...)

		ts, err := template.New(name).Funcs(functions).ParseFiles(files...)
		if err != nil {
			return myCache, err
		}

		myCache[name] = ts
	}

	return myCache, nil
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

func Iterate(count int) []int {
	var items []int
	for i := 1; i <= count; i++ {
		items = append(items, i)
	}
	return items
}

func Add(a, b int) int {
	return a + b
}
