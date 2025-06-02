# About

Experimenting Kafka in Go

# How to run

## Set up Kafka brokers and utilities

Start docker compose

```
docker compose up
```

## Consumer

Start docker compose

Start app

```
go run ./cmd/consumer/main.go
```

## Producer

Start docker compose

Start app

```
go run ./cmd/producer/main.go
```

# Experiments

## Experiment 1

What happens if consumer continuously reads message from topic(s) and there is another goroutine to make that consumer subscribe different topic(s)

Result: After resubscribng to new topic(s), consumer finishes the current message of the first topic(s) (if have one) and change to read message from the new topic(s)

## Experiment 2

How much is difference of performance between using equal amount of consumers and less amount of consumers compared to amount of topic partitions

Result: -
