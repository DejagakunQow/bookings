package repository

import (
	"time"

	"github.com/DejagakunQow/bookings/cmd/web/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)

	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	AllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)

	CountReservations() (int, error)
	CountNewReservations() (int, error)

	GetReservationByID(id int) (models.Reservation, error)
	GetReservationsForMonth(start, end time.Time) ([]models.Reservation, error)

	UpdateReservation(u models.Reservation) error
	DeleteReservation(id int) error
	AllReservationsForDate(date time.Time) ([]models.Reservation, error)
	UpdateProcessedForReservation(id, processed int) error
	IsRoomAvailable(roomID int, start, end time.Time) (bool, error)
	AllRooms() ([]models.Room, error)
}
