package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	servers := "localhost:29092"
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
	})
	if err != nil {
		log.Fatalf("Error while creating producer: %s", err)
		return
	}
	defer producer.Close()
	log.Printf("Created producer")

	// Create topics
	createTopic("test3", servers, 5)
	createTopic("test4", servers, 1)

	// Produce messages
	count := 500_000
	produceNMessages(producer, "test3", count)
	produceNMessages(producer, "test4", count)
}

func createTopic(topic string, servers string, partitions int) {
	// Create a topic with given parition number
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": servers,
	})
	if err != nil {
		log.Fatalf("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}
	defer adminClient.Close()
	specifications := []kafka.TopicSpecification{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = adminClient.CreateTopics(
		ctx,
		specifications,
		kafka.SetAdminOperationTimeout(5*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create topic: %s\n", err)
	} else {
		log.Printf("Created topic [%s] with [%d] partitions", topic, partitions)
	}
}

func produceNMessages(producer *kafka.Producer, topic string, count int) {
	// Produce messages
	log.Printf("Starting producing %d messages to topic %s", count, topic)
	for range count {
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("ABCDE"),
		}, nil)
		if err != nil {
			log.Fatalf("Error while producing message: %s", err)
		}
	}

	// Flush outstanding messages
	// - Ensures all messages in the producer's internal queue are sent
	// - Ensures delivery reports are processed before the program exits
	producer.Flush(10_000)
	log.Printf("Produced %d messages to topic %s", count, topic)
}
