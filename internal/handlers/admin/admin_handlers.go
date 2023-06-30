package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/itzYubi/bookings/internal/config"
	"github.com/itzYubi/bookings/internal/driver"
	"github.com/itzYubi/bookings/internal/forms"
	"github.com/itzYubi/bookings/internal/helpers"
	"github.com/itzYubi/bookings/internal/models"
	"github.com/itzYubi/bookings/internal/render"
	"github.com/itzYubi/bookings/internal/repository"
	"github.com/itzYubi/bookings/internal/repository/dbrepo"
)

// AdminRepo the repository used by the admin handlers
var AdminRepo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewAdminRepo creates a new repository
func NewAdminRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewAdminHandlers(r *Repository) {
	AdminRepo = r
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservations shows all new reservations for admin
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservations gets all reservations for admin
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation shows a particular reservation
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	//get reservation from database
	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "admin-show-reservation.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

// AdminPostShowReservation Updates the reservation in DB
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println("Cannot parse the form")
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Phone = r.Form.Get("phone")
	res.Email = r.Form.Get("email")

	err = m.DB.UpdateReservation(res)
	if err != nil {
		m.App.ErrorLog.Println("Updating the reservation in DB failed")
		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Saved Changes")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminProcessReservation marks reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "src")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		m.App.ErrorLog.Println("Updating the Processed value in DB failed")
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminDeleteReservation deletes a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	src := chi.URLParam(r, "src")
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.DB.DeleteReservation(id)
	if err != nil {
		m.App.ErrorLog.Println("Delete reservation with id: " + strconv.Itoa(id) + " in DB FAILED")
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation Deleted")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminPostReservationsCalendar handles POST of reservation calendar
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println("Cannot parse the form")
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	rooms, err := m.DB.AllRooms()
	if err != nil {
		m.App.ErrorLog.Println("Cannot get all rooms from DB")
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, room := range rooms {

		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		for name, value := range curMap {
			if value > 0 {
				if !form.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, name)) {
					m.App.InfoLog.Println("Delete the block: " + name + ", with value: " + strconv.Itoa(value))
					err := m.DB.DeleteRoomBlockByID(value)
					if err != nil {
						m.App.ErrorLog.Println("Deletion of block: " + name + ", with value: " + strconv.Itoa(value) + "failed in DB with error: " + err.Error())
					}
				}
			}
		}
	}

	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomId, _ := strconv.Atoi(exploded[2])
			dateString := exploded[3]
			date, _ := time.Parse("2006-01-2", dateString)
			m.App.InfoLog.Println("Insert block for room Id: " + strconv.Itoa(roomId) + ", for date: " + date.String())
			err := m.DB.InsertBlockForRoom(roomId, date)
			if err != nil {
				m.App.ErrorLog.Println("Insertion of block: " + name + ", failed in DB with error: " + err.Error())
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes Saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}

// AdminReservationsCalendar displays the  resservation calendar for admin
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	//assume there is no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastYear

	stringMap["this_month"] = fmt.Sprintf("%02d", now.Month())
	stringMap["this_month_name"] = now.Month().String()
	stringMap["this_month_year"] = now.Format("2006")

	currentYear, currentMonth, _ := now.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	//get all rooms from DB
	rooms, err := m.DB.AllRooms()
	if err != nil {
		m.App.ErrorLog.Println("Cannot get all rooms data from DB: rooms")
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	//get restriction for rooms using maps
	isFirst := true
	dayMap := make(map[int]string)

	for _, room := range rooms {
		//create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0

			if isFirst {
				dayMap[d.Day()] = d.Weekday().String()
			}
		}
		isFirst = false

		//get all restrictions for current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(room.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			m.App.ErrorLog.Println("Getting all restrictions by date failed in the DB. Error: " + err.Error())
			helpers.ServerError(w, err)
			return
		}

		for _, restriction := range restrictions {
			if restriction.ReservationID > 0 {
				//It is a reservation
				for d := restriction.StartDate; !d.After(restriction.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = restriction.ReservationID
				}
			} else {
				//It is a block
				for d := restriction.StartDate; !d.After(restriction.EndDate); d = d.AddDate(0, 0, 1) {
					blockMap[d.Format("2006-01-2")] = restriction.ID
				}
			}
		}
		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)
	}

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}
