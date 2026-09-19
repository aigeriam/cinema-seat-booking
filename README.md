# Cinema Seat Booking

A small Go service for holding and booking cinema seats. Redis is used as the shared store so multiple requests can safely compete for the same seat.

## Run

```bash
docker compose up -d redis
go run .
```

The application runs at `http://localhost:6969`. Redis Commander is available at `http://localhost:8081` when the full Compose stack is running.

## API

- `GET /movies` lists the available movies.
- `GET /movies/{movieID}/seats` lists seats currently held or confirmed for a movie.
- `POST /movies/{movieID}/seats/{seatID}/hold` holds a seat. Body: `{"user_id":"user-1"}`.
- `PUT /sessions/{sessionID}/confirm` confirms a hold. Body: `{"user_id":"user-1"}`.
- `DELETE /sessions/{sessionID}` releases a hold. Body: `{"user_id":"user-1"}`.

## Redis implementation

Each seat is stored under `seat:{movieID}:{seatID}`. A hold is created with Redis `SET NX`. Holds expire after two minutes.

The service also stores `session:{sessionID}` pointing to the seat key. This lets confirmation and release find the seat from the session ID.

- **Hold:** creates the seat key and session key with a two-minute TTL.
- **Confirm:** removes both TTLs with `PERSIST` and changes the booking status to `confirmed`.
- **Release:** deletes both keys so the seat becomes available again.

The booking service depends on a `BookingStore` interface, so the Redis store can be replaced with another implementation for tests or different deployments.