package models

import (
	"time"
)

// User is the user model
type User struct {
	ID         int
	FirstName  string
	LastName   string
	Email      string
	Password   string
	AccesLevel int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Room is the Rooms Model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restrictions Model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is the reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// RoomRestrictions is the room restriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	RestrictionID int
	ReservationID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	Reseration    Reservation
	Restriction   Restriction
}

// MailData holds email message
type MailData struct {
	To       string
	From     string
	Subject  string
	Content  string
	Template string
}
