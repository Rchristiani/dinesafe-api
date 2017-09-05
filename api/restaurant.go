package dinesafe

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const connectionString string = "user=ryanchristiani dbname=dinesafe sslmode=disable"

//Query represents the data in the Dinesafe dataset
type Query struct {
	Rows []Row `xml:"ROW"`
}

//Row is for each restaurant
type Row struct {
	RowID                     int    `xml:"ROW_ID" json:"rowID"`
	EstablishmentID           int    `xml:"ESTABLISHMENT_ID" json:"establishmentID"`
	InspectionID              int    `xml:"INSPECTION_ID" json:"inspectionID"`
	EstablishmentName         string `xml:"ESTABLISHMENT_NAME" json:"establishmentName"`
	EstablishmentType         string `xml:"ESTABLISHMENTTYPE" json:"establishmentType"`
	EstablishmentAddress      string `xml:"ESTABLISHMENT_ADDRESS" json:"establishmentAddress"`
	EstablishmentStatus       string `xml:"ESTABLISHMENT_STATUS" json:"establishmentStatus"`
	MinimumInspectionsPerYear int    `xml:"MINIMUM_INSPECTIONS_PERYEAR" json:"MinimumInspectionsPerYear"`
	InfractionDetails         string `xml:"INFRACTION_DETAILS" json:"infractionsDetails"`
	InspectionDate            string `xml:"INSPECTION_DATE" json:"inpectionDate"`
	Severity                  string `xml:"SEVERITY" json:"severity"`
	Action                    string `xml:"ACTION" json:"action"`
	CourtOutcome              string `xml:"COURT_OUTCOME" json:"courtOutcome"`
	AmountFinded              string `xml:"AMOUNT_FINED" json:"amountFinded"`
}

func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length")
}

//GetRestaurants gets all the restaurants
func GetRestaurants(w http.ResponseWriter, r *http.Request) {

	SetHeaders(w)

	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	if offset == "" {
		offset = "0"
	}

	if limit == "" {
		limit = "50"
	}

	db, err := sql.Open("postgres", connectionString)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM restaurants LIMIT $1 OFFSET $2;", limit, offset)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}
	var restaurants []Row
	for rows.Next() {
		var restaurant Row
		err = rows.Scan(
			&restaurant.RowID,
			&restaurant.EstablishmentID,
			&restaurant.InspectionID,
			&restaurant.EstablishmentName,
			&restaurant.EstablishmentType,
			&restaurant.EstablishmentAddress,
			&restaurant.EstablishmentStatus,
			&restaurant.MinimumInspectionsPerYear,
			&restaurant.InfractionDetails,
			&restaurant.InspectionDate,
			&restaurant.Severity,
			&restaurant.Action,
			&restaurant.CourtOutcome,
			&restaurant.AmountFinded,
		)

		if err != nil {
			log.Fatal(err)
		}

		restaurants = append(restaurants, restaurant)
	}
	restaurantJSON, err := json.Marshal(restaurants)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(restaurantJSON)

}

func GetRestaurantByID(w http.ResponseWriter, r *http.Request) {

	SetHeaders(w)

	vars := mux.Vars(r)

	id := vars["id"]

	db, err := sql.Open("postgres", connectionString)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT * FROM restaurants WHERE id=$1", id)

	var restaurant Row

	err = row.Scan(
		&restaurant.RowID,
		&restaurant.EstablishmentID,
		&restaurant.InspectionID,
		&restaurant.EstablishmentName,
		&restaurant.EstablishmentType,
		&restaurant.EstablishmentAddress,
		&restaurant.EstablishmentStatus,
		&restaurant.MinimumInspectionsPerYear,
		&restaurant.InfractionDetails,
		&restaurant.InspectionDate,
		&restaurant.Severity,
		&restaurant.Action,
		&restaurant.CourtOutcome,
		&restaurant.AmountFinded,
	)

	if err != nil {
		log.Fatal(err)
	}

	restoJSON, err := json.Marshal(restaurant)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(restoJSON)
}
