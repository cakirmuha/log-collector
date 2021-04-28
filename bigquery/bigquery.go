package bigquery

import (
	"cloud.google.com/go/bigquery"
	"context"
	"log"
	"os"
)

type BigQuery struct {
	ProjectId string
	DatasetId string
	TableId   string
	Client    *bigquery.Client
}

func NewBigQuery() *BigQuery {
	ctx := context.Background()

	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Println("There is no project id set in env variable")
		return nil
	}

	datasetId := os.Getenv("DATASET_ID")
	if datasetId == "" {
		log.Println("There is no dataset id set in env variable")
		return nil
	}

	tableId := os.Getenv("TABLE_ID")
	if tableId == "" {
		log.Println("There is no table id set in env variable")
		return nil
	}

	client, err := bigquery.NewClient(ctx, projectId)
	if err != nil {
		log.Println("Could not create bigquery Client: %v", err)
		return nil
	}

	return &BigQuery{
		ProjectId: projectId,
		DatasetId: datasetId,
		Client:    client,
	}
}
