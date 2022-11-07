package main

import (
	"github.com/itzYubi/bookings/config"
	"github.com/itzYubi/bookings/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	//mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	//mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	return mux
}
