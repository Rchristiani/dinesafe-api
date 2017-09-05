package main

import (
	"database/sql"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	dinesafe "github.com/rchristiani/dinesafe/api"
)

func load() {
	//Open the XML file
	dinesafeXML, err := os.Open("dinesafe.xml")
	//Make sure to defer and close he files
	defer dinesafeXML.Close()
	if err != nil {
		log.Fatal(err)
	}
	//Take the file and read the bytes from it
	xmlBytes, err := ioutil.ReadAll(dinesafeXML)
	if err != nil {
		log.Fatal(err)
	}
	//Make a struct for everything
	var rows dinesafe.Query
	//Unmarshal the bytes into the Query
	xml.Unmarshal(xmlBytes, &rows)

	db, err := sql.Open("postgres", "user=ryanchristiani dbname=dinesafe sslmode=disable")

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, restaurant := range rows.Rows {
		_, err = db.Exec("INSERT INTO restaurants(establishmentID, inspectionID, establishmentName, establishmentType, establishmentAddress, establishmentStatus, minimumInspectionsPerYear, infractionDetails, inspectionDate, severity, action, courtOutcome, amountFinded) Values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13);",
			restaurant.EstablishmentID, restaurant.InspectionID, restaurant.EstablishmentName, restaurant.EstablishmentType, restaurant.EstablishmentAddress, restaurant.EstablishmentStatus, restaurant.MinimumInspectionsPerYear, restaurant.InfractionDetails, restaurant.InspectionDate, restaurant.Severity, restaurant.Action, restaurant.CourtOutcome, restaurant.AmountFinded,
		)

		if err != nil {
			log.Fatalln(err)
		}

	}

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/restaurants", dinesafe.GetRestaurants).Methods("GET")
	r.HandleFunc("/restaurants/{id}", dinesafe.GetRestaurantByID).Methods("GET")

	http.ListenAndServe(":3700", r)
}
