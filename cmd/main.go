package main

import (
	"cinema-seat-booking/internal/adapters/redis"
	"cinema-seat-booking/internal/booking"
	handler "cinema-seat-booking/internal/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fmt.Println("http://localhost:6969")
	mux.HandleFunc("GET /movies", handler.ListMovies)
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	client := redis.NewClient("localhost:6379")
	store := booking.NewRedisStore(client)
	svc := booking.NewService(store)
	bookingHandler := handler.NewHandler(svc)
	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeats)
	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)
	if err := http.ListenAndServe(":6969", mux); err != nil {
		log.Fatal(err)
	}
}
