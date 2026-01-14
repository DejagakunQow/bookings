package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"encoding/gob"

	"github.com/DejagakunQow/bookings/cmd/web/internal/config"
	"github.com/DejagakunQow/bookings/cmd/web/internal/driver"
	"github.com/DejagakunQow/bookings/cmd/web/internal/handlers"
	"github.com/DejagakunQow/bookings/cmd/web/internal/helpers"
	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
	"github.com/DejagakunQow/bookings/cmd/web/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// production flag
	app.InProduction = false

	// session setup
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	helpers.NewHelpers(&app)

	// -----------------------------
	// DATABASE CONNECTION (FIX)
	// -----------------------------
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := driver.ConnectSQL(dbURL)

	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	defer db.SQL.Close()

	// -----------------------------
	// TEMPLATE CACHE
	// -----------------------------
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	// -----------------------------
	// REPOSITORIES & RENDERER
	// -----------------------------

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	// -----------------------------
	// HTTP SERVER
	// -----------------------------
	fmt.Printf("Starting application on port %s\n", port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() (*driver.DB, error) {

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbname := flag.String("dbname", "", "Database name")
	dbuser := flag.String("dbuser", "", "Database user")
	dbhost := flag.String("dbhost", "localhost", "Database host")
	dbport := flag.String("dbport", "5432", "Database port")
	dbpass := flag.String("dbpass", "", "Database password")
	dbSSL := flag.String("dbssl", "disable", "Database SSL settings(disable, prefer, require)")

	flag.Parse()

	if *dbname == "" || *dbuser == "" {
		fmt.Println("Missing required flags - dbname and dbuser")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbhost, *dbport, *dbname, *dbuser, *dbpass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil

}
