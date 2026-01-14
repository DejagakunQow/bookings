package main

import (
	"net/http"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// Middleware
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad(app))

	// Static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Public routes
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/choose-room", handlers.Repo.ChooseRoom)
	mux.Get("/contact", handlers.Repo.Contact)

	// Auth routes
	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// Reservation routes
	mux.Get("/make-reservation", handlers.Repo.Reservation)

	// Admin routes
	mux.Route("/admin", func(mux chi.Router) {
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)

		// 🧠 Calendar route — this makes the page accessible
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)

		mux.Get("/reservations-new", handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservations)

		mux.Get("/process-reservations/{src}/{id}", handlers.Repo.AdminProcessReservation)

		mux.Get("/reservation/{src}/{id}", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)
		mux.Get("/reservations/all/{id}", handlers.Repo.AdminShowReservation)

		mux.Get("/reservations-delete/{id}", handlers.Repo.AdminDeleteReservation)
		mux.Post("/calendar/update/{id}", handlers.Repo.UpdateCalendarReservation)
		mux.Post("/calendar/save", handlers.Repo.UpdateCalendarReservation)
		mux.Post("/calendar/delete/{id}", handlers.Repo.DeleteCalendarReservation)
	})

	return mux
}
