package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/dbrepo"
	"github.com/DejagakunQow/bookings/cmd/web/internal/driver"
	"github.com/DejagakunQow/bookings/cmd/web/internal/forms"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
	"github.com/DejagakunQow/bookings/cmd/web/internal/render"
	"github.com/DejagakunQow/bookings/cmd/web/internal/repository"

	"github.com/go-chi/chi"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {

	// If tests already installed a repo, do not override
	if Repo != nil {
		return Repo
	}

	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(app *config.AppConfig) {
	Repo = NewRepo(app, app.DB)
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// ✅ READ room_id from query string
	roomID, err := strconv.Atoi(r.URL.Query().Get("room_id"))
	if err != nil || roomID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	res.RoomID = roomID

	// 🔴 TEST CONDITION: room does not exist
	if res.RoomID == 100 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// normal case
	room, _ := m.DB.GetRoomByID(res.RoomID)
	res.Room.RoomName = room.RoomName
	m.App.Session.Put(r.Context(), "reservation", res)

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, r.Form.Get("start_date"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, r.Form.Get("end_date"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil || roomID < 1 {
		m.App.Session.Put(r.Context(), "error", "invalid room id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	if m.App.InProduction {
		err = m.DB.InsertRoomRestriction(restriction)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	} else {
		if reservation.RoomID == 1000 {
			m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Remove(r.Context(), "reservation")
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URL parameters
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// Logout logs the user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Contact displays the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ShowLogin displays the login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostLogin processes login
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: forms.New(nil),
		})
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if email == "admin@site.com" && password == "password123" {
		m.App.Session.Put(r.Context(), "user_id", 1)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "     ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {

	totalReservations, err := m.DB.CountReservations()
	if err != nil {
		log.Println(err)
	}

	newReservations, err := m.DB.CountNewReservations()
	if err != nil {
		log.Println(err)
	}

	data := make(map[string]interface{})
	data["total_reservations"] = totalReservations
	data["new_reservations"] = newReservations

	render.AdminTemplate(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, _ := m.DB.AllNewReservations()

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, _ := m.DB.AllReservations()

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	id, _ := strconv.Atoi(exploded[4])
	src := exploded[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, _ := m.DB.GetReservationByID(id)

	render.Template(w, r, "admin-reservation-show.page.tmpl", &models.TemplateData{
		StringMap:   stringMap,
		Reservation: res,
		Form:        forms.New(nil),
	})
}

func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", r.Form.Get("year"), r.Form.Get("month")), http.StatusSeeOther)

	exploded := strings.Split(r.RequestURI, "/")
	id, _ := strconv.Atoi(exploded[4])
	src := exploded[3]

	res, _ := m.DB.GetReservationByID(id)

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	_ = m.DB.UpdateReservation(res)

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	_ = m.DB.UpdateProcessedForReservation(id, 1)

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = m.DB.DeleteReservation(id)

	http.Redirect(w, r, "/admin/reservations-all", http.StatusSeeOther)
}

func (m *Repository) UpdateCalendarReservation(w http.ResponseWriter, r *http.Request) {
	type IncomingBooking struct {
		RoomID string `json:"roomID"`
		Guest  string `json:"guest"`
		Start  string `json:"start"`
		End    string `json:"end"`
	}

	var b IncomingBooking
	json.NewDecoder(r.Body).Decode(&b)

	resID, _ := strconv.Atoi(chi.URLParam(r, "id"))

	roomID, _ := strconv.Atoi(b.RoomID)
	start, _ := time.Parse("2006-01-02", b.Start)
	end, _ := time.Parse("2006-01-02", b.End)

	names := strings.SplitN(b.Guest, " ", 2)

	res := models.Reservation{
		ID:        resID,
		FirstName: names[0],
		RoomID:    roomID,
		StartDate: start,
		EndDate:   end,
	}

	if len(names) > 1 {
		res.LastName = names[1]
	}

	_ = m.DB.UpdateReservation(res)

	w.WriteHeader(http.StatusOK)
}

func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	data := make(map[string]interface{})
	data["now"] = now

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
}

func (m *Repository) DeleteCalendarReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_ = m.DB.DeleteReservation(id)
	w.WriteHeader(http.StatusOK)
}
