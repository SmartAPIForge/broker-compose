package schema_uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"
)

var schemas = map[string]string{
	"NewZip":        path.Join("avro", "newzip.avsc"),
	"ProjectStatus": path.Join("avro", "project-status.avsc"),
	"DeployPayload": path.Join("avro", "deploy-payload.avsc"),
}

type SchemaUploader struct {
	mu                sync.RWMutex
	schemaRegistryURL string
}

func NewSchemaUploader(schemaRegistryURL string) *SchemaUploader {
	manager := &SchemaUploader{
		schemaRegistryURL: schemaRegistryURL,
	}

	return manager
}

func (sm *SchemaUploader) UploadSchemasToRegistry() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for topic, schemaPath := range schemas {
		if err := sm.uploadSchemaToRegistry(topic, schemaPath); err != nil {
			panic(fmt.Sprintf("Failed to upload schema for topic %s: %v", topic, err))
		}
	}
}

func (sm *SchemaUploader) uploadSchemaToRegistry(topic, schemaPath string) error {
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	type schemaRequest struct {
		Schema string `json:"schema"`
	}
	requestBody, err := json.Marshal(schemaRequest{Schema: string(schemaData)})
	if err != nil {
		return err
	}

	schemaURL := fmt.Sprintf("%s/subjects/%s-value/versions", sm.schemaRegistryURL, topic)

	resp, err := http.Get(schemaURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		fmt.Printf("Schema for topic %s already exists in registry\n", topic)
		return nil
	}

	req, err := http.NewRequest("POST", schemaURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.schemaregistry.v1+json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return err
	}

	fmt.Printf("Schema for topic %s successfully uploaded to registry\n", topic)
	return nil
}
