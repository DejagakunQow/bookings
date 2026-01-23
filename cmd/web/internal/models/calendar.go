package models

import "time"

type CalendarDay struct {
	Day          int
	Date         time.Time
	Reservations []Reservation
}
