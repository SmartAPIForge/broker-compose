package main

import (
	"broker-compose/src/schema-uploader"
	"broker-compose/src/topic-creator"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	schemaRegistryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	broker := os.Getenv("KAFKA_HOST")

	uploader := schema_uploader.NewSchemaUploader(schemaRegistryURL)
	uploader.UploadSchemasToRegistry()

	creator := topic_creator.NewTopicCreator(broker)
	creator.CreateTopics()
}
