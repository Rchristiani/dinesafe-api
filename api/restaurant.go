package dinesafe

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Restaurants represents the data in the Dinesafe dataset
type Restaurants struct {
	Rows []Row `xml:"ROW"`
}

//Row is for each restaurant
type Row struct {
	RowID                     int    `xml:"ROW_ID" json:"rowID"`
	EstablishmentID           int    `xml:"ESTABLISHMENT_ID" json:"establishmentID"`
	EstablishmentName         string `xml:"ESTABLISHMENT_NAME" json:"establishmentName"`
	EstablishmentType         string `xml:"ESTABLISHMENTTYPE" json:"establishmentType"`
	EstablishmentAddress      string `xml:"ESTABLISHMENT_ADDRESS" json:"establishmentAddress"`
	EstablishmentStatus       string `xml:"ESTABLISHMENT_STATUS" json:"establishmentStatus"`
	MinimumInspectionsPerYear int    `xml:"MINIMUM_INSPECTIONS_PERYEAR" json:"minimumInspectionsPerYear"`
}

func scanRestaurantRows(res *Row, rows *sql.Rows) error {
	err := rows.Scan(
		&res.RowID,
		&res.EstablishmentID,
		&res.EstablishmentName,
		&res.EstablishmentType,
		&res.EstablishmentAddress,
		&res.EstablishmentStatus,
		&res.MinimumInspectionsPerYear,
	)

	return err
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

	db, err := sql.Open("postgres", ConnectionString)
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
		err = scanRestaurantRows(&restaurant, rows)

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

//GetRestaurantByID is used to get a restaurant by ID
func GetRestaurantByID(w http.ResponseWriter, r *http.Request) {

	SetHeaders(w)

	vars := mux.Vars(r)

	id := vars["id"]

	db, err := sql.Open("postgres", ConnectionString)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT * FROM restaurants WHERE id=$1", id)

	var restaurant Row

	err = row.Scan(
		&restaurant.RowID,
		&restaurant.EstablishmentID,
		&restaurant.EstablishmentName,
		&restaurant.EstablishmentType,
		&restaurant.EstablishmentAddress,
		&restaurant.EstablishmentStatus,
		&restaurant.MinimumInspectionsPerYear,
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

//SearchRestaurantsByName takes a name query and returns matching restaurants
func SearchRestaurantsByName(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	name := r.URL.Query().Get("name")

	db, err := sql.Open("postgres", ConnectionString)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	query := fmt.Sprintf("SELECT * FROM Restaurants WHERE establishmentName LIKE '%%%s%%';", name)

	rows, err := db.Query(query)

	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var restaurants []Row

	for rows.Next() {
		var restaurant Row

		err = scanRestaurantRows(&restaurant, rows)

		restaurants = append(restaurants, restaurant)

		if err != nil {
			log.Fatal(err)
		}
	}

	restaurantsJSON, err := json.Marshal(restaurants)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(restaurantsJSON)
}
