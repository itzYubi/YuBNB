package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/itzYubi/bookings/internal/config"
	"github.com/itzYubi/bookings/internal/driver"
	"github.com/itzYubi/bookings/internal/handlers"
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

	fmt.Printf("Staring application on port %s", webPort)

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

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false

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
	log.Println("Connecting to Database.....")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=root")
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
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
