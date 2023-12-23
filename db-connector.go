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

type Project struct {
	Es  *opensearch.Client
}

var project = Project{}

func init() {
	// Initialize and connect to OpenSearch
	config := opensearch.Config{
		Addresses: []string{
			os.Getenv("OPENSEARCH_URL"),
		},
		Username: os.Getenv("OPENSEARCH_USERNAME"),
		Password: os.Getenv("OPENSEARCH_PASSWORD"),
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

	log.Printf("[INFO] Saving task %s to Opensearch", task.ID)

	err = indexEs(ctx, "Tasks", task.ID, body)
	if err != nil {
		log.Printf("[ERROR] Error saving task %s to Opensearch: %s", task.ID, err)
		return err
	}

	return nil
}

func SaveTaskResult(ctx context.Context, result Result) error {
	log.Printf("Saving task result for task %s to Opensearch", result.TaskID)

	// convert to string
	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("[ERROR] Error marshalling task result for task %s: %s", result.TaskID, err)
		return err
	}

	err = indexEs(ctx, "Results", result.TaskID, resultBytes)
	if err != nil {
		return err
	}

	return nil
}

func GetTaskResult(ctx context.Context, taskID string) (Result, error) {
	result := Result{}

	log.Printf("Getting task result for task %s from Opensearch", taskID)
	res, err := project.Es.Get(strings.ToLower("results"), taskID)
	if err != nil {
		return result, err
	}

	if res.IsError() {
		log.Printf("[ERROR] Error getting task result for task %s from Opensearch: %s", taskID, res.String())
		return result, errors.New(res.String())
	}

	if res.StatusCode != 200 {
		return result, errors.New(fmt.Sprintf("Bad statuscode from database: %d for task %s", res.StatusCode, taskID))
	}

	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("[WARNING] Error reading the response body from Opensearch for task %s: %s", taskID, err)
		return result, err
	}

	log.Printf("[INFO] Got openserach response for task %s: %s", taskID, string(respBody))

	resultWrapper := ResultWrapper{}

	err = json.Unmarshal(respBody, &resultWrapper)
	if err != nil {
		log.Printf("[ERROR] Error parsing the response body from Opensearch for task %s: %s. Raw: %s", taskID, err, respBody)
		return result, err
	}

	result = resultWrapper.Source

	return result, nil
}
