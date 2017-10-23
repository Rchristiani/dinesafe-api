package dinesafe

import "net/http"

const ConnectionString string = "user=ryanchristiani dbname=dinesafe sslmode=disable"

//SetHeaders sets headers for responses
func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length")
}
