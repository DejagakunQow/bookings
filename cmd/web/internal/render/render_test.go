package render

import (
	"net/http"
	"os"
	"testing"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
	"github.com/alexedwards/scs/v2"
)

var testApp config.AppConfig

func TestMain(m *testing.M) {
	testApp.InProduction = false

	session := scs.New()
	testApp.Session = session

	App = &testApp

	os.Exit(m.Run())
}

func TestAddDefaultData(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	// 🔑 LOAD SESSION INTO CONTEXT (THIS IS THE FIX)
	ctx, err := App.Session.Load(req.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)

	// now session works
	App.Session.Put(ctx, "flash", "123")

	td := &models.TemplateData{}

	result := AddDefaultData(td, req)

	if result.Flash != "123" {
		t.Errorf("expected flash value of 123, got %s", result.Flash)
	}
}
