package main

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
	})
	if err != nil {
		log.Fatalf("Error while creating producer: %s", err)
		return
	}
	defer producer.Close()
	log.Printf("Created producer")

	// Produce messages
	count := 1000
	produceNMessages(producer, "test1", count)
	produceNMessages(producer, "test2", count)
}

func produceNMessages(producer *kafka.Producer, topic string, count int) {
	log.Printf("Starting producing %d messages to topic %s", count, topic)
	for range count {
		err := producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(topic),
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
