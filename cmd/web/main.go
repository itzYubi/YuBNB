package main

import (
	"fmt"
	"github.com/itzYubi/bookings/config"
	"github.com/itzYubi/bookings/handlers"
	"github.com/itzYubi/bookings/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const webPort = ":7070"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {
	// change this to true when in production
	app.InProduction = false

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

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
