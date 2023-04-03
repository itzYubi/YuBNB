package main

import (
	"net/http"

	"github.com/itzYubi/bookings/internal/config"
	"github.com/itzYubi/bookings/internal/handlers"

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
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{Rid}", handlers.Repo.ChooseRoom)

	mux.Get("/contact", handlers.Repo.Contact)

	fileserver := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))
	return mux
}
