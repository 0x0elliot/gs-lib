package gslib

import (
	// "context"
	"encoding/json"

	"github.com/opensearch-project/opensearch-go"
	// "github.com/opensearch-project/opensearch-go/opensearchapi"
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

func saveTask(ctx, task Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = client.Index(indexName).DocumentID(task.ID).BodyJson(string(body)).Do(ctx)
	return err
}

func getTaskResult(ctx, taskID string) (*Result, error) {
	res, err := client.Get().Index(indexName).DocumentID(taskID).Do(ctx)
	if err != nil {
		return nil, err
	}

	var result Result
	err = json.Unmarshal([]byte(res.Body), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
