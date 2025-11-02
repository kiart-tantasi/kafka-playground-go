# About

Experimenting Kafka in Go

# Start Kafka docker compose

Start docker compose

```sh
docker compose up
```

# Experiments

## Experiment 1

What happens if consumer continuously reads message from topic(s) and there is another goroutine to make that consumer subscribe different topic(s)

### Steps

- Produce 1000 messages into topic `test1` and topic `test2`

```sh
go run ./cmd/producer/main.go
```

- Consume messages while switching between topic `test1` and `test2` every 500 ms.

```sh
go run ./cmd/consumer/main.go
```

### Result

After resubscribing to new topic(s), consumer finishes the current message of the first topic(s) (if have one) and change to read message from the new topic(s)

## Experiment 2

How much is difference of performance between 1 paritions and 5 partitions with 5 consumers

### Steps

- Create topics with x partitions and produce messages into them

```sh
go run ./cmd/producer2/main.go
```

- Consume messages from topic that has 5 partitions

```sh
go run ./cmd/consumer2/main.go test3
```

- Consume messages from topic that has 1 partition

```sh
go run ./cmd/consumer2/main.go test4
```

### Result

5 paritions have 28% faster performance than 1 parition. The experiment is done with 500,000 messages and each message contains 5-character string"".
