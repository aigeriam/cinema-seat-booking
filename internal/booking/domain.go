package booking

import "errors"

type Booking struct {
	ID      string
	UserID  string
	MovieID string
	SeatID  string
	Status  string
}

// for dependency injection
type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieID string) []Booking
}

var (
	ErrSeatAlreadyBooked = errors.New("The seat is already booked")
)
