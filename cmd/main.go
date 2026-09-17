package cmd

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
	mux.HandleFunc("GET /movies", ListMovies)
	
}
func ListMovies(w http.ResponseWriter, r *http.Request){
	WriteJSON
}