package booking

import (
	"context"
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
	Book(b Booking) (Booking, error)
	ListBookings(movieID string) []Booking
	Confirm(ctx context.Context, sessionId string, userID string) (Booking, error)
	Release(ctx context.Context, sessionId string, userID string) error
}

var (
	ErrSeatAlreadyBooked      = errors.New("The seat is already booked")
	ErrSeatWasntHoldToConfirm = errors.New("The seat wasnt hold to be confirmed")
)
