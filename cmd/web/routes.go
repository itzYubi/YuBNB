package main

import (
	"net/http"

	"github.com/itzYubi/bookings/config"
	"github.com/itzYubi/bookings/handlers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/natsu-home", handlers.Repo.NatsuHome)
	mux.Get("/yuki-home", handlers.Repo.YukiHome)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Get("/contact", handlers.Repo.Contact)


	fileserver := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))
	return mux
}
