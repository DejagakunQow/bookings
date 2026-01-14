package models

import "github.com/DejagakunQow/bookings/cmd/web/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float64
	Data            map[string]interface{}
	Reservation     Reservation
	IsAuthenticated bool
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
}
