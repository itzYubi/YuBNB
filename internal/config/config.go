package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/itzYubi/bookings/internal/models"
)

type AppConfig struct {
	InProduction  bool
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
