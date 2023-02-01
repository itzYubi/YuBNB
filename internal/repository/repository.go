package repository

import "github.com/itzYubi/bookings/internal/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error) 
	InsertRoomRestriction(r models.RoomRestriction) error
}
