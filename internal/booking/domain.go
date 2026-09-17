package booking

import (
	"errors"
	"time"
)

type Booking struct {
	ID        string
	UserID    string
	MovieID   string
	SeatID    string
	Status    string
	ExpiresAt time.Time
}

// for dependency injection
type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}

var (
	ErrSeatAlreadyBooked = errors.New("The seat is already booked")
)
