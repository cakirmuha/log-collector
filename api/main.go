package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	lc "github.com/cakirmuha/log-collector"
	"github.com/cakirmuha/log-collector/bigquery"
	"github.com/cakirmuha/log-collector/pub_sub"
	"github.com/gorilla/mux"
)

var (
	pubSub   *pub_sub.PubSub
	bigQuery *bigquery.BigQuery
)

func createNewLog(w http.ResponseWriter, r *http.Request) {
	var (
		res        lc.ApiResponse
		codewayLog lc.CodewayLog
		msgData    []byte
	)

	if err := func() error {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(reqBody, &codewayLog)
		if err != nil {
			return err
		}

		eventTime := codewayLog.EventTime
		if eventTime == nil {
			return fmt.Errorf("event_time is nil")
		}
		tsString := fmt.Sprintf("%.3f", float64(*eventTime)/float64(1000))
		codewayLog.EventTime = nil
		codewayLog.EventTs = &tsString

		msgData, err = json.Marshal(codewayLog)
		if err != nil {
			return err
		}
		defer func() {
			err = r.Body.Close()
			return
		}()

		return err
	}(); err != nil {
		log.Printf("Parse request body: %v\n", err.Error())
		errorMessage := "Parse request body"

		res = lc.ApiResponse{
			Success: false,
			Error:   &errorMessage,
		}
		_ = json.NewEncoder(w).Encode(res)
	}

	if err := pubSub.Publish(msgData); err != nil {
		log.Printf("Publish message: %v\n", err.Error())
		errorMessage := "Publish message to pub/sub"
		res = lc.ApiResponse{
			Success: false,
			Error:   &errorMessage,
		}
		_ = json.NewEncoder(w).Encode(res)
		return
	}
	res = lc.ApiResponse{
		Success: true,
		Error:   nil,
	}
	_ = json.NewEncoder(w).Encode(res)
}

func getAnalytics(w http.ResponseWriter, r *http.Request) {
	var (
		res       lc.ApiResponse
		analytics lc.CodewayAnalytics
	)
	dailyActiveUsers, err := bigQuery.DailyActiveUsers()
	if err != nil {
		log.Printf("DailyActiveUsers: %v\n", err.Error())
		errorMessage := "DailyActiveUsers calculation error"

		res = lc.ApiResponse{
			Success: false,
			Error:   &errorMessage,
		}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	dailyAverageDurations, err := bigQuery.DailyAverageDurations()
	if err != nil {
		log.Printf("DailyAverageDurations: %v\n", err.Error())
		errorMessage := "DailyAverageDurations calculation error"

		res = lc.ApiResponse{
			Success: false,
			Error:   &errorMessage,
		}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	totalUsers, err := bigQuery.TotalUsers()
	if err != nil {
		log.Printf("TotalUsers: %v\n", err.Error())
		errorMessage := "TotalUsers calculation error"

		res = lc.ApiResponse{
			Success: false,
			Error:   &errorMessage,
		}
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	if totalUsers != nil {
		analytics.TotalUser = *totalUsers
	}
	if dailyActiveUsers != nil {
		analytics.DailyActiveUsers = dailyActiveUsers
	}
	if dailyAverageDurations != nil {
		analytics.DailyAverageDurations = dailyAverageDurations
	}

	res = lc.ApiResponse{
		Success: true,
		Error:   nil,
		Data:    analytics,
	}

	_ = json.NewEncoder(w).Encode(res)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/log", createNewLog).Methods("POST")
	myRouter.HandleFunc("/analytics", getAnalytics).Methods("GET")
	log.Fatal(http.ListenAndServe(":8484", myRouter))
}

func main() {
	pubSub = pub_sub.NewPubSub()
	if pubSub == nil {
		log.Fatal("Error creating Pub/Sub client")
		return
	}
	bigQuery = bigquery.NewBigQuery()
	if bigQuery == nil {
		log.Fatal("Error creating BigQuery client")
		return
	}
	handleRequests()
}
