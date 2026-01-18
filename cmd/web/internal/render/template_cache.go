package render

import (
	"html/template"
	"path/filepath"
)

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	isoDate := 0
	funcMap := template.FuncMap{
		"humanDate": humanDate,
		"isoDate":   isoDate,
	}

	// Page templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return cache, err
	}

	// Public layouts
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

		ts, err := template.New(name).Funcs(funcMap).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseFiles(layouts...)
			if err != nil {
				return cache, err
			}
		}

		if len(adminLayouts) > 0 {
			ts, err = ts.ParseFiles(adminLayouts...)
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
