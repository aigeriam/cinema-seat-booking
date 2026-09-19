package booking

import (
	"cinema-seat-booking/internal/booking"
	"cinema-seat-booking/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type handler struct {
	svc *booking.Service
}
type Request struct {
	UserID string `json:"user_id"`
}

func NewHandler(svc *booking.Service) *handler {
	return &handler{svc}
}

type sessionResponse struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	SeatID    string `json:"seat_id"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

func (h *handler) ListSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	bookings := h.svc.ListBookings(movieID)
	seats := make([]seatInfo, 0, len(bookings))
	for _, b := range bookings {
		seats = append(seats, seatInfo{
			SeatID: b.SeatID,
			UserID: b.UserID,
			Booked: true,
		})
	}
	utils.WriteJSON(w, http.StatusOK, seats)
}

func (h *handler) HoldSeats(w http.ResponseWriter, r *http.Request) {
	movieID := r.PathValue("movieID")
	seatID := r.PathValue("seatID")
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	data := booking.Booking{
		UserID:  req.UserID,
		SeatID:  seatID,
		MovieID: movieID,
	}
	session, err := h.svc.Book(data)
	if err != nil {
		log.Println(err)
		return
	}
	type holdResponse struct {
		SessionID string `json:"session_id"`
		MovieID   string `json:"movie_id"`
		SeatID    string `json:"seat_id"`
		ExpiresAt string `json:"expires_at"`
	}
	utils.WriteJSON(w, http.StatusCreated, holdResponse{
		SeatID:    seatID,
		MovieID:   session.MovieID,
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	})

}
func (h *handler) ConfirmSession(w http.ResponseWriter, r *http.Request) {

	sessionid := r.PathValue("sessionID")
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	if req.UserID == "" {
		return
	}
	session, err := h.svc.Confirm(r.Context(), sessionid, req.UserID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, "oops")
		return
	}

	utils.WriteJSON(w, http.StatusAccepted, sessionResponse{
		SessionID: session.ID,
		MovieID:   session.MovieID,
		SeatID:    session.SeatID,
		UserID:    req.UserID,
		Status:    session.Status,
	})
}

func (h *handler) ReleaseSession(w http.ResponseWriter, r *http.Request) {
	sessionid := r.PathValue("sessionID")
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		return
	}
	if req.UserID == "" {
		return
	}
	err := h.svc.Release(r.Context(), sessionid, req.UserID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
