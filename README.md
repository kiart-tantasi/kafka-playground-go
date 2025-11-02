# About

Experimenting Kafka in Go

# How to run

## Set up Kafka brokers and utilities

Start docker compose

```sh
docker compose up
```

## Producer

### Produce 1000 messages into topic `test1` and topic `test2`

```sh
go run ./cmd/producer/main.go
```

## Consumer

### Experiment 1 - Consume messages while switching between topic `test1` and `test2` every 500 ms.

```sh
go run ./cmd/consumer/main.go
```

### Experiment 2 - 5 Consumers consume messages from 1 paritions vs 5 paritions

(in-progress)

```sh
go run ./cmd/consumer2/main.go
```

# Experiments

## Experiment 1

What happens if consumer continuously reads message from topic(s) and there is another goroutine to make that consumer subscribe different topic(s)

Result: After resubscribing to new topic(s), consumer finishes the current message of the first topic(s) (if have one) and change to read message from the new topic(s)

## Experiment 2

How much is difference of performance between using equal amount of consumers and less amount of consumers compared to amount of topic partitions

Result: -
