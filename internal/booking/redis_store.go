package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"uuid"

	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute //time for selected seat to be hold

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}

}
func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}
func (s *RedisStore) Book(b Booking) (Booking, error) {
	session, err := s.hold(b)
	if err != nil {
		return Booking{}, err
	}
	log.Printf("Session booked %v", session)
	return session, nil
}
func (s *RedisStore) ListBookings(MovieID string) []Booking {
	pattern := fmt.Sprintf("seat:%s:*", MovieID)
	var sessions []Booking
	ctx := context.Background()
	iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		val, err := s.rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}
		session, err := parseSession(val)
		if err != nil {
			continue
		}
		sessions = append(sessions, session)
	}
	return sessions
}
func (s *RedisStore) hold(b Booking) (Booking, error) {
	id := uuid.New().String()
	now := time.Now()
	ctx := context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)
	b.ID = id
	val, _ := json.Marshal(b)
	res := s.rdb.SetArgs(ctx, key, val, redis.SetArgs{
		Mode: "NX", // set, if not exists
		TTL:  defaultHoldTTL,
	})
	ok := res.Val() == "OK"
	if !ok {
		return Booking{}, ErrSeatAlreadyBooked
	}
	//entry for session
	s.rdb.Set(ctx, sessionKey(id), key, defaultHoldTTL)
	return Booking{
		ID:        id,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		UserID:    b.UserID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}, nil
}

// Persist removes the ttl, so session never expires
func (s *RedisStore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	booking, seatKey, err := s.getsession(ctx, sessionID, userID)
	if err != nil {
		return Booking{}, err
	}
	s.rdb.Persist(ctx, sessionKey(sessionID))
	s.rdb.Persist(ctx, seatKey)
	booking.Status = "confirmed"
	data := Booking{
		ID:      booking.ID,
		MovieID: string(booking.MovieID),
		SeatID:  booking.SeatID,
		UserID:  booking.UserID,
		Status:  "confirmed",
	}
	js, _ := json.Marshal(data)
	s.rdb.Set(ctx, seatKey, js, 0)
	return booking, nil

}

func (s *RedisStore) getsession(ctx context.Context, sessionID string, userID string) (Booking, string, error) {
	seatKey, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return Booking{}, "", err
	}
	val, err := s.rdb.Get(ctx, seatKey).Result()
	if err != nil {
		return Booking{}, "", err
	}
	booking, err := parseSession(val)
	if err != nil {
		return Booking{}, "", err
	}
	return booking, seatKey, nil

}
func (s *RedisStore) Release(ctx context.Context, sessionID string, userID string) error {
	_, seatKey, err := s.getsession(ctx, sessionID, userID)
	if err != nil {
		return err
	}
	s.rdb.Del(ctx, seatKey, sessionKey(sessionID))
	return nil
}
func parseSession(val string) (Booking, error) {
	var data Booking
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return Booking{}, err
	}
	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil
}
