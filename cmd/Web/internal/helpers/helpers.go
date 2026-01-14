package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	if app != nil && app.ErrorLog != nil {
		app.ErrorLog.Println(trace)
	} else {
		// fallback: never panic
		fmt.Println(trace)
	}

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
func HumanDate(t time.Time) string {
	return t.Format("02 Jan 2006")
}
