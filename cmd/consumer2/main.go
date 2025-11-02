package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// 1 partition - 9232 ms
// 5 paritions - 7207 ms

func main() {

	// Get topic argument
	topic := ""
	if len(os.Args) < 1 {
		log.Fatal("No topic argument defined. Please run with argument(s) `<topic>`")
	} else {
		topic = os.Args[1]
		log.Printf("App will start with topic %s in 3 seconds", topic)
		time.Sleep(3 * time.Second)
	}

	// Configs
	groupId := fmt.Sprintf("group-id-random-%d", time.Now().UnixMilli())

	// Log configs
	log.Printf("Group id [%s], app will start in 3 seconds..", groupId)
	time.Sleep(3 * time.Second)

	// Create 5 consumers
	consumers := []*kafka.Consumer{}
	for i := range 5 {
		log.Printf("Starting consumer indexed %d", i)

		// Create a consumer and let it subscribe a topic
		consumer, err := createConsumer(groupId, "localhost:29092")
		if err != nil {
			log.Fatalf("Error while creating consumer: %v", err)
		}
		err = consumer.Subscribe(topic, nil)
		if err != nil {
			log.Fatalf("Error while subscribing a topic: %v", err)
		}
		consumers = append(consumers, consumer)
	}

	// Counters
	var counter int32 = 0
	start := time.Now()

	// Start consuming in 5 goroutines
	for i := range 5 {
		go func(i int) {
			// Retrieve a consumer
			consumer := consumers[i]
			defer consumer.Close()

			for {
				_, err := consumer.ReadMessage(-1)
				if err != nil {
					log.Printf("Error while reading message: %s", err)
					continue
				} else {
					atomic.AddInt32(&counter, 1)

					// log every x messages
					if counter%100 == 0 {
						log.Printf("counter: %d, taken %d ms", counter, time.Since(start).Milliseconds())
					}
				}
			}
		}(i)
	}

	for {
		log.Println("Consumers are waiting for messages")
		time.Sleep(10 * time.Second)
	}

}

func createConsumer(groupId, servers string) (*kafka.Consumer, error) {
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
}
