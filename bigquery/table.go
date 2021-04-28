package bigquery

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// ImportClusteredTable demonstrates creating a table and define partitioning and clustering properties.
func (bq *BigQuery) CreateEventTable() error {
	var err error

	// define daily ingestion-based partitioning.
	dayPartitioning := &bigquery.TimePartitioning{
		Type:  bigquery.DayPartitioningType,
		Field: "event_ts",
	}

	clustering := &bigquery.Clustering{Fields: []string{
		"session_id", "event_ts",
	}}

	schema := bigquery.Schema{
		{Name: "event_ts", Type: bigquery.TimestampFieldType},
		{Name: "session_id", Type: bigquery.StringFieldType},
		{Name: "user_id", Type: bigquery.StringFieldType},
		{Name: "type", Type: bigquery.StringFieldType},
		{Name: "app_id", Type: bigquery.StringFieldType},
		{Name: "event_name", Type: bigquery.StringFieldType},
		{Name: "page", Type: bigquery.StringFieldType},
		{Name: "country", Type: bigquery.StringFieldType},
		{Name: "region", Type: bigquery.StringFieldType},
		{Name: "city", Type: bigquery.StringFieldType},
	}

	if err := bq.Client.Dataset(bq.DatasetId).Table(bq.TableId).Create(context.Background(), &bigquery.TableMetadata{
		Schema:           schema,
		TimePartitioning: dayPartitioning,
		Clustering:       clustering,
	}); err != nil {
		return fmt.Errorf("create biquery table with error: %v", err.Error())
	}

	return err
}
