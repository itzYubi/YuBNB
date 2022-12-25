package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/itzYubi/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Errorf("unexpected type: %T, expected *chi.Mux", v)
	}
}
