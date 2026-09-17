package booking

type Service struct {
	store BookingStore ///bookingstore is interface
}

func NewService(store BookingStore) *Service {
	return &Service{store}
}

func (s *Service) Book(b Booking) error {
	return s.store.Book(b)
}
