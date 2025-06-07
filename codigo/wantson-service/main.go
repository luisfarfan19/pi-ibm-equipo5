package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	handler "wantson-service/internal/services"
)

func main() {
	//err := db.ConnectMongo()
	//if err != nil {
	//	return
	//}

	r := mux.NewRouter()
	r.HandleFunc("/health", handler.HealthHandler).Methods("GET")
	r.HandleFunc("/watson/validate/planogram", handler.ValidatePlanogramHandler).Methods("POST")
	r.HandleFunc("/watson/validate/shelve", handler.ValidateShelterHandler).Methods("POST")
	r.HandleFunc("/store/planogram-response/section", handler.GetSectionResponseHandler).Methods("GET")
	r.HandleFunc("/watson/getToken", handler.GetIbmAccessTokenHandler).Methods("GET")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/stores", handler.GetStoresHandler).Methods("GET")
	r.HandleFunc("/store/data", handler.GetStoreDataHandler).Methods("GET")

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
