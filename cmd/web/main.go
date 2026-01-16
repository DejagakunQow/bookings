package main

import (
	"database/sql"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/driver"
	"github.com/DejagakunQow/bookings/cmd/web/internal/handlers"
	"github.com/DejagakunQow/bookings/cmd/web/internal/helpers"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
	"github.com/DejagakunQow/bookings/cmd/web/internal/render"
)

var app config.AppConfig

func main() {

	// ------------------------------------------------
	// Register types for session serialization
	// ------------------------------------------------
	registerGOB()

	// ------------------------------------------------
	// Load environment variables
	// ------------------------------------------------
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// ------------------------------------------------
	// Connect to database
	// ------------------------------------------------
	db, err := connectDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	// ------------------------------------------------
	// App configuration
	// ------------------------------------------------
	app.InProduction = os.Getenv("GO_ENV") == "production"
	app.DB = db

	// ------------------------------------------------
	// Session manager (CRITICAL SECTION)
	// ------------------------------------------------
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// 🔴 REQUIRED: attach a store, otherwise app WILL panic
	session.Store = postgresstore.New(db.SQL)

	app.Session = session

	// ------------------------------------------------
	// Initialize helpers, renderers, handlers
	// ------------------------------------------------
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	handlers.NewHandlers(&app)

	// ------------------------------------------------
	// HTTP server
	// ------------------------------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on port", port)
	log.Println("Initializing routes")

	handler := routes(&app)

	log.Println("Routes initialized, starting HTTP server")

	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}

// ------------------------------------------------
// Helper functions
// ------------------------------------------------

func registerGOB() {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})
}

func connectDB(dsn string) (*driver.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &driver.DB{SQL: db}, nil
}
