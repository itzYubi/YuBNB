package main

import (
	"net/http"

	"github.com/itzYubi/bookings/internal/config"
	"github.com/itzYubi/bookings/internal/handlers"
	adminHandlers "github.com/itzYubi/bookings/internal/handlers/admin"

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
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact", handlers.Repo.Contact)
	mux.Post("/submit-contact", handlers.Repo.PostSubmitContact)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	fileserver := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", adminHandlers.AdminRepo.AdminDashboard)
		mux.Get("/reservations-new", adminHandlers.AdminRepo.AdminNewReservations)
		mux.Get("/reservations-all", adminHandlers.AdminRepo.AdminAllReservations)
		mux.Get("/reservations-calendar", adminHandlers.AdminRepo.AdminReservationsCalendar)
		mux.Post("/reservations-calendar", adminHandlers.AdminRepo.AdminPostReservationsCalendar)
		mux.Get("/process-reservation/{src}/{id}/process", adminHandlers.AdminRepo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/delete", adminHandlers.AdminRepo.AdminDeleteReservation)

		mux.Get("/reservations/{src}/{id}/show", adminHandlers.AdminRepo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", adminHandlers.AdminRepo.AdminPostShowReservation)
	})
	return mux
}
