package scripts

import (
	"github.com/cakirmuha/log-collector/bigquery"
	"log"
)

func createBigqueryTable() {
	bigQuery := bigquery.NewBigQuery()
	if bigQuery == nil {
		log.Fatal("Error creating BigQuery client")
		return
	}

	err := bigQuery.CreateEventTable()
	if err != nil {
		log.Printf("Error creating eventtable %v\n", err.Error())
	}
	return
}
