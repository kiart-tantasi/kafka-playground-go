package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Configs
	var INTERVAL = 500 * time.Millisecond
	groupId := fmt.Sprintf("group-id-random-%d", time.Now().UnixMilli())

	// Log configs
	log.Printf("Group id [%s] with interval [%s], app will start in 3 seconds..", groupId, INTERVAL)
	time.Sleep(3 * time.Second)

	// Create consumer
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Error while creating consumer: %s", err)
		return
	}
	defer consumer.Close()
	log.Printf("Created consumer with group id %s", groupId)

	// Subcribe initial topics
	initialTopics := []string{"test1"}
	consumer.SubscribeTopics(initialTopics, nil)
	log.Printf("Subscribed initial topic(s) %s", initialTopics)

	// Topic switcher
	go func() {
		for {
			time.Sleep(INTERVAL)

			var currentTopics []string
			if currentTopics, err = consumer.Subscription(); err != nil {
				log.Printf("Error when getting subscription: %s", err)
				currentTopics = []string{}
			}

			// switch
			if currentTopics[0] == "test1" {
				currentTopics[0] = "test2"
			} else {
				currentTopics[0] = "test1"
			}

			consumer.SubscribeTopics(currentTopics, nil)
			// log.Printf("Resubscribed topic(s) %s", currentTopics)
		}
	}()

	topicTest1 := 0
	topicTest2 := 0

	for {
		// Consume a message
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error while reading message: %s", err)
			continue
		}

		// Process message
		message := string(msg.Value)
		// Fake process time
		time.Sleep(1 * time.Millisecond)

		// Count
		// log.Printf("Message: %s", message)
		if message == "test1" {
			topicTest1++
		} else if message == "test2" {
			topicTest2++
		} else {
			log.Fatalf("Found unexpcted messages [%s]", message)
		}
		log.Printf("topicTest1: [%d], topicTest2: [%d]", topicTest1, topicTest2)
	}

}
