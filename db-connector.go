// db-connector.go
package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

const (
	indexName = "seo_tasks"
)

var client *opensearch.Client

func init() {
	// Initialize and connect to OpenSearch
	client, _ = opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://opensearch:9200"},
	})
}

func saveTask(task Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = client.Index(indexName).DocumentID(task.ID).BodyJson(string(body)).Do(context.Background())
	return err
}

func getTaskResult(taskID string) (*TaskResult, error) {
	res, err := client.Get().Index(indexName).DocumentID(taskID).Do(context.Background())
	if err != nil {
		return nil, err
	}

	var result TaskResult
	err = json.Unmarshal([]byte(res.Body), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
