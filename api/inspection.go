package dinesafe

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Inspections struct {
	Inspections []Inspection `xml:"ROW"`
}

type Inspection struct {
	EstablishmentID   int    `xml:"ESTABLISHMENT_ID" json:"establishmentID"`
	InspectionID      int    `xml:"INSPECTION_ID" json:"inspectionID"`
	InfractionDetails string `xml:"INFRACTION_DETAILS" json:"infractionsDetails"`
	InspectionDate    string `xml:"INSPECTION_DATE" json:"inpectionDate"`
	Severity          string `xml:"SEVERITY" json:"severity"`
	Action            string `xml:"ACTION" json:"action"`
	CourtOutcome      string `xml:"COURT_OUTCOME" json:"courtOutcome"`
	AmountFinded      string `xml:"AMOUNT_FINED" json:"amountFinded"`
}

func scanInspectionRows(ins *Inspection, rows *sql.Rows) error {
	err := rows.Scan(
		&ins.EstablishmentID,
		&ins.InspectionID,
		&ins.InfractionDetails,
		&ins.InspectionDate,
		&ins.Severity,
		&ins.Action,
		&ins.CourtOutcome,
		&ins.AmountFinded,
	)

	return err
}

//GetInspections gets inpsections
func GetInspections(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT * FROM Inspections LIMIT $1 OFFSET $2", limit, offset)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var inspections []Inspection

	for rows.Next() {
		var inspection Inspection

		err := scanInspectionRows(&inspection, rows)

		if err != nil {
			log.Fatal(err)
		}

		inspections = append(inspections, inspection)
	}

	JSON, err := json.Marshal(inspections)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(JSON)

}

//GetInspectionByResID will get inspections for a specific restraunt
func GetInspectionByResID(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	vars := mux.Vars(r)

	id := vars["id"]

	db, err := sql.Open("postgres", ConnectionString)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM Inspections WHERE EstablishmentID=$1", id)
	defer rows.Close()

	if err != nil {
		log.Fatal(err)
	}

	var inspections []Inspection

	for rows.Next() {
		var inspection Inspection

		err = scanInspectionRows(&inspection, rows)

		if err != nil {
			log.Fatal(err)
		}

		inspections = append(inspections, inspection)

	}

	JSON, err := json.Marshal(inspections)

	if err != nil {
		log.Fatal(err)
	}

	w.Write(JSON)

}
