package topic_creator

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var topics = []string{"NewZip", "ProjectStatus"}

type TopicCreator struct {
	broker string
}

func NewTopicCreator(broker string) *TopicCreator {
	return &TopicCreator{
		broker: broker,
	}
}

func (tc *TopicCreator) CreateTopics() {
	adminConfig := kafka.ConfigMap{
		"bootstrap.servers": tc.broker,
	}

	adminClient, err := kafka.NewAdminClient(&adminConfig)
	if err != nil {
		panic(fmt.Sprintf("error creating Kafka AdminClient: %v", err))
	}
	defer adminClient.Close()

	var topicSpecs []kafka.TopicSpecification
	for _, topic := range topics {
		topicSpecs = append(topicSpecs, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	results, err := adminClient.CreateTopics(context.Background(), topicSpecs)
	if err != nil {
		panic(fmt.Errorf("error creating topics: %v", err))
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			log.Printf("Failed creating topic %s: %v", result.Topic, result.Error)
			continue
		}
		log.Printf("Successfully created topic %s", result.Topic)
	}
}
