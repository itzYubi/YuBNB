package dbrepo

import (
	"errors"
	"time"

	"github.com/itzYubi/bookings/internal/models"
)

// InsertReservation inserts a reservation to the reservation database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 0 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == -1 {
		return errors.New("some error")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for roomID, and false if no availability exists.
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomById gets the room by id
func (m *testDBRepo) GetRoomById(roomId int) (models.Room, error) {
	var room models.Room
	if roomId > 2 {
		return room, errors.New("some error")
	}
	return room, nil
}
