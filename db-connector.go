package gslib

import (
	"os"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"crypto/tls"
	"net/http"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

const (
	indexName = "seo_tasks"
)

type Project struct {
	Es  *opensearch.Client
}

var project = Project{}

func init() {
	// Initialize and connect to OpenSearch
	config := opensearch.Config{
		Addresses: []string{
			os.Getenv("OPENSEARCH_URL",	"https://localhost:9200"),
		},
		Username: os.Getenv("OPENSEARCH_USERNAME", "admin"),
		Password: os.Getenv("OPENSEARCH_PASSWORD", "admin"),
	}

	transport := http.DefaultTransport.(*http.Transport)
	transport.MaxIdleConnsPerHost = 100
	transport.ResponseHeaderTimeout = time.Second * 10
	transport.Proxy = nil


	transport.TLSClientConfig = &tls.Config{
		MinVersion:         tls.VersionTLS11,
		InsecureSkipVerify: true,
	}

	config.Transport = transport

	client, err := opensearch.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating the opensearch client: %s", err)
	}

	project.Es = client
}

func indexEs(ctx context.Context, nameKey string, id string, bytes []byte) error {
	req := opensearchapi.IndexRequest{
		Index:      strings.ToLower(nameKey),
		DocumentID: id,
		Body:       strings.NewReader(string(bytes)),
		Refresh:    "true",
		Pretty:     true,
	}

	res, err := req.Do(ctx, project.Es)
	if err != nil {
		return err
	}

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
		return errors.New(res.String())
	}

	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[WARNING] Error reading the response body from Opensearch: %s", err)
		return err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 {
		return errors.New(fmt.Sprintf("Bad statuscode from database: %d. Reason: %s", res.StatusCode, string(respBody)))
	}

	var r map[string]interface{}
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		log.Printf("[WARNING] Error parsing the response body from Opensearch: %s. Raw: %s", err, respBody)
		return err
	}

	return nil
}

func SaveTask(ctx context.Context, task Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	log.Printf("Saving task %s to Opensearch", task.ID)

	err = indexEs(ctx, indexName, task.ID, body)
	if err != nil {
		return err
	}

	return nil
}

func GetTaskResult(ctx context.Context, taskID string) (*Result, error) {
	log.Printf("Getting task result for task %s from Opensearch", taskID)
	res, err := project.Es.Get(strings.ToLower(indexName), taskID)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[WARNING] Error reading the response body from Opensearch for task %s: %s", taskID, err)
		return nil, err
	}

	var result Result
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
