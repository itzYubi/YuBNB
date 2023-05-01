package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hako/durafmt"
	"github.com/itzYubi/bookings/internal/config"
	"github.com/itzYubi/bookings/internal/driver"
	"github.com/itzYubi/bookings/internal/forms"
	"github.com/itzYubi/bookings/internal/helpers"
	"github.com/itzYubi/bookings/internal/models"
	"github.com/itzYubi/bookings/internal/render"
	"github.com/itzYubi/bookings/internal/repository"
	"github.com/itzYubi/bookings/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find room by id!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	duration := (res.EndDate).Sub(res.StartDate)
	durStr, _ := durafmt.ParseString(duration.String())
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	stringMap["duration"] = durStr.String()

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert reservation into DB!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	//insert room restriction
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't insert restriction into DB!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	//send mails
	dateFormat := "2006-01-02"
	htmlGuestMessage := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s:, <br>
		This is to confirm your reservartion at WeebzHome Bookings, with the following details:<br>
		Home name: %s,<br>
		Arrival: %s,<br>
		Departure: %s,<br>
		<br>
		Thanks & regards,<br>
		Yubi,<br>
		(WeebZ admin)
	`, reservation.FirstName, reservation.Room.RoomName, reservation.StartDate.Format(dateFormat), reservation.EndDate.Format(dateFormat))

	htmlOwnerMessage := fmt.Sprintf(`
		<strong>Reservation Notification</strong><br>
		Dear owner:, <br>
		This is to notify that you have a reservartion from WeebzHome Bookings, with the following details:<br>
		Home name: %s,<br>
		Arrival: %s,<br>
		Departure: %s,<br>
		<br>
		Thanks & regards,<br>
		Yubi,<br>
		(WeebZ admin)
	`, reservation.Room.RoomName, reservation.StartDate.Format(dateFormat), reservation.EndDate.Format(dateFormat))

	//send notifications - first to guest
	msg := models.MailData{
		To:       reservation.Email,
		From:     "bookingsAdmin@Weebz.com",
		Subject:  "Reservation Confirmation for: " + reservation.FirstName + " " + reservation.LastName,
		Content:  htmlGuestMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	//send notifications - second to owner

	msg = models.MailData{
		To:       "roomOwner@world.com",
		From:     "bookingsAdmin@Weebz.com",
		Subject:  "Reservation Notification for: " + reservation.Room.RoomName,
		Content:  htmlOwnerMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) NatsuHome(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "natsu.page.tmpl", &models.TemplateData{})
}

func (m *Repository) YukiHome(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "yuki.page.tmpl", &models.TemplateData{})
}

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostSearchAvailability searches availabilty of all rooms and displayes available room in next screen
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	for _, room := range rooms {
		m.App.InfoLog.Println("ROOM:", room.ID, room.RoomName)
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK         bool   `json:"ok"`
	Message    string `json:"message"`
	RoomId     string `json:"room_id"`
	Start_date string `json:"start_date"`
	End_date   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and writes a JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	jr := jsonResponse{
		OK:         available,
		RoomId:     strconv.Itoa(roomID),
		Start_date: sd,
		End_date:   ed,
		Message:    "",
	}

	out, err := json.MarshalIndent(jr, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Println(string(out))

	w.Header().Add("Content-Type", "application/json")
	w.Write(out)
}

// displays the reservation summary
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "Rid"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("failed to get reservation from session"))
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes URl parameters, builds a sessional variable, and takes user to make reservation screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	var res models.Reservation

	res.RoomID = roomID
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := m.DB.GetRoomById(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// Contact renders the contact page.
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	var contactData models.ContactData
	data := make(map[string]interface{})
	data["contact"] = contactData

	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) PostSubmitContact(w http.ResponseWriter, r *http.Request) {
	var contactData models.ContactData
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	contactData.Name = r.Form.Get("name")
	contactData.Email = r.Form.Get("email")
	contactData.Subject = r.Form.Get("subject")
	contactData.Message = r.Form.Get("message")

	form := forms.New(r.PostForm)

	form.Required("name", "email", "subject", "message")
	form.MinLength("name", 5)
	form.MinLength("message", 10)
	form.IsEmail("email")

	if !form.IsValid() {
		data := make(map[string]interface{})
		data["contact"] = contactData
		err := render.Template(w, r, "contact.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		log.Println(err)
		return
	}

	htmlSupportMessage := fmt.Sprintf(`
		<strong>Help Request</strong><br>
		Dear Support Team, <br>
		This is support request by customer, with the following details:<br>
		Name: %s,<br>
		Email: %s,<br>
		Support Message: %s,<br>
		<br>
		Thanks,<br>
		(WeebZ admin)
	`, contactData.Name, contactData.Email, contactData.Message)

	msg := models.MailData{
		To:       "support.weebz@gmail.com",
		From:     "contact.weebz@gmail.com",
		Subject:  contactData.Subject,
		Content:  htmlSupportMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg
	m.App.Session.Put(r.Context(), "flash", "We have received your contact request!")

	http.Redirect(w, r, "/contact", http.StatusSeeOther)
}
