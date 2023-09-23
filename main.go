package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/orders", GetOrders).Methods("GET")
	r.HandleFunc("/orders/{id}/update-status", UpdateShipmentStatus).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
