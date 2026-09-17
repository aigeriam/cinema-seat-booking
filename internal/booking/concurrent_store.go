package booking

import "sync"

type ConcurrentStore struct {
	sync.RWMutex
	bookings map[string]Booking ///the key is seatid
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{bookings: map[string]Booking{}}
}

// /writing
func (s *ConcurrentStore) Book(b Booking) error {
	s.Lock() //the next user has to wait untill we unlock
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}
	s.bookings[b.SeatID] = b
	return nil
}
func (s ConcurrentStore) ListBookings(movieID string) []Booking {
	s.RLock()
	defer s.RUnlock() //other go routines can read the data

	var result []Booking
	for _, b := range s.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}
	return result
}
