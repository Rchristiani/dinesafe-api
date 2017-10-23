package main

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	dinesafe "github.com/rchristiani/dinesafe/api"
)

func main() {
	r := mux.NewRouter()

	//Restraunts
	r.HandleFunc("/restaurants", dinesafe.GetRestaurants).Methods("GET")
	r.HandleFunc("/restaurants/search", dinesafe.SearchRestaurantsByName).Methods("GET")
	r.HandleFunc("/restaurants/{id}", dinesafe.GetRestaurantByID).Methods("GET")

	//Inspections
	r.HandleFunc("/inspections", dinesafe.GetInspections).Methods("GET")
	r.HandleFunc("/inspections/{id}", dinesafe.GetInspectionByResID).Methods("GET")

	http.ListenAndServe(":3700", r)
}
