package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	handler "wantson-service/internal/services"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", handler.HealthHandler).Methods("GET")
	r.HandleFunc("/watson/validate/shelve", handler.ValidateShelterHandler).Methods("GET")

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
