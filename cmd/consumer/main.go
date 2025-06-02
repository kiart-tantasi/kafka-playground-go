package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Create consumer
	groupId := fmt.Sprintf("group-id-random-%d", time.Now().UnixMilli())
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
			time.Sleep(500 * time.Millisecond)

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
			log.Printf("Resubscribed topic(s) %s", currentTopics)
		}
	}()

	counterTest1 := 0
	counterTest2 := 0
	counterOther := 0

	for {
		// Info
		// if subcribedTopics, err := consumer.Subscription(); err != nil {
		// 	log.Printf("Error while getting subscription info: %s", err)
		// 	continue
		// } else {
		// 	log.Printf("Reading a message from topics %s", subcribedTopics)
		// }

		// Read
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
			counterTest1++
		} else if message == "test2" {
			counterTest2++
		} else {
			counterOther++
		}
		log.Printf("counterTest1: [%d], counterTest2: [%d], counterOther: [%d]", counterTest1, counterTest2, counterOther)
	}

}
