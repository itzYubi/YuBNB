package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/itzYubi/bookings/internal/config"
	"github.com/itzYubi/bookings/internal/driver"
	"github.com/itzYubi/bookings/internal/handlers"
	adminHandlers "github.com/itzYubi/bookings/internal/handlers/admin"
	"github.com/itzYubi/bookings/internal/helpers"
	"github.com/itzYubi/bookings/internal/models"
	"github.com/itzYubi/bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const webPort = ":7070"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	db, err := run()
	if err != nil {
		log.Println(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	log.Println("Staring mail listener...")
	listenForMail()

	log.Printf("Staring application on port %s \n", webPort)

	srv := &http.Server{
		Addr:    webPort,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	//things to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.ContactData{})
	gob.Register(map[string]int{})

	//read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database Password")
	dbPort := flag.String("dbport", "5432", "Database Port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings(disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required Flags")
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "\nINFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "\nERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to DB
	log.Println("Connecting to Database....")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to Database! Dying...")
	}

	log.Println("CONNECTED TO DATABASE!")
	//------------------------------------------------------------------------

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache due to " + err.Error())
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	adminrepo := adminHandlers.NewAdminRepo(&app, db)
	handlers.NewHandlers(repo)
	adminHandlers.NewAdminHandlers(adminrepo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
