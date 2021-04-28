package bigquery

import (
	"context"
	"log"

	lc "github.com/cakirmuha/log-collector"
	"google.golang.org/api/iterator"
)

func (bq *BigQuery) DailyActiveUsers() ([]lc.DailyActiveUser, error) {
	query := bq.Client.Query("SELECT DATE(event_ts) as dte, COUNT(DISTINCT user_id) as cnt " +
		"FROM event " +
		"GROUP BY dte " +
		"ORDER BY dte")
	query.DefaultProjectID = bq.ProjectId
	query.DefaultDatasetID = bq.DatasetId
	it, err := query.Read(context.Background())
	if err != nil {
		log.Printf("Error on reading query results: %v", err)
		return nil, err
	}

	var dailyActiveUsers []lc.DailyActiveUser
	for {
		var d lc.DailyActiveUser
		err := it.Next(&d)
		if err == iterator.Done {
			break
		}
		d.Day = d.Dte.String()
		if err != nil {
			log.Printf("Error on iterating query results: %v", err)
			return nil, err
		}
		dailyActiveUsers = append(dailyActiveUsers, d)
	}

	return dailyActiveUsers, err
}

func (bq *BigQuery) DailyAverageDurations() ([]lc.DailyAverageDuration, error) {
	query := bq.Client.Query("SELECT s.dte, SUM(s.diff) as duration FROM (" +
		"SELECT DATE(event_ts) as dte, session_id, TIMESTAMP_DIFF(max(event_ts), min(event_ts), SECOND) as diff " +
		"FROM event " +
		"GROUP BY dte, session_id " +
		"ORDER BY dte, session_id) AS s GROUP BY s.dte ORDER BY s.dte")

	query.DefaultProjectID = bq.ProjectId
	query.DefaultDatasetID = bq.DatasetId
	it, err := query.Read(context.Background())
	if err != nil {
		log.Printf("Error on reading query results: %v", err)
		return nil, err
	}

	var dailyAverageDurations []lc.DailyAverageDuration
	for {
		var d lc.DailyAverageDuration
		err := it.Next(&d)
		if err == iterator.Done {
			break
		}
		d.Day = d.Dte.String()
		if err != nil {
			log.Printf("Error on iterating query results: %v", err)
			return nil, err
		}
		dailyAverageDurations = append(dailyAverageDurations, d)
	}

	return dailyAverageDurations, err
}

func (bq *BigQuery) TotalUsers() (*lc.TotalUser, error) {
	query := bq.Client.Query("SELECT COUNT(DISTINCT user_id) as cnt " +
		"FROM event ")
	query.DefaultProjectID = bq.ProjectId
	query.DefaultDatasetID = bq.DatasetId
	it, err := query.Read(context.Background())
	if err != nil {
		log.Printf("Error on reading query results: %v", err)
		return nil, err
	}

	var totalUser *lc.TotalUser
	for {
		var t lc.TotalUser
		err := it.Next(&t)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error on iterating query results: %v", err)
			return nil, err
		}
		totalUser = &t
	}

	return totalUser, nil
}
