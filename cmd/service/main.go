package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"kafka-parallel-queues/pkg/handler"
	"kafka-parallel-queues/pkg/kafkakit"

	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

func readQueuesNumberFromEnv() int {
	queuesString := os.Getenv("QUEUES")
	if queuesString == "" {
		return 1
	}

	queues, err := strconv.Atoi(queuesString)
	if err != nil {
		return 1
	}

	return queues
}

func readConsumersNumberFromEnv() int {
	consumersString := os.Getenv("CONSUMERS")
	if consumersString == "" {
		return 1
	}

	consumers, err := strconv.Atoi(consumersString)
	if err != nil {
		return 1
	}

	return consumers
}

func main() {
	k := os.Getenv("KAFKA_HOST")
	fmt.Println("service/main, listening ", k)

	queues := readQueuesNumberFromEnv()
	consumers := readConsumersNumberFromEnv()

	consumerGroup, ctx := errgroup.WithContext(context.Background())
	fetchContext, stopReading := context.WithCancel(ctx)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit

		slog.Info("stop reading messages")

		stopReading()
	}()

	for i := 0; i < consumers; i++ {
		consumerGroup.Go(func() error {
			runConsumer(fetchContext, queues)

			return nil
		})
	}

	slog.Info("service is running")

	if err := consumerGroup.Wait(); err != nil {
		slog.Error("failed to stop reading by consumer group", "error", err)

		return
	}

	slog.Info("service successfully stopped")
}

func runConsumer(fetchContext context.Context, queues int) {
	const queueSize = 3

	handlerPool := handler.NewPool(queues, queueSize)

	defer func() {
		if err := handlerPool.Close(); err != nil {
			slog.Error("failed to close handler pool", "error", err)
		}
	}()

	k := os.Getenv("KAFKA_HOST")
	reader := kafkakit.NewReader(fetchContext, kafka.ReaderConfig{
		Brokers: []string{k},
		GroupID: "numbers-consumer-group-0",
		Topic:   "numbers",
	})

	defer func() {
		if err := reader.Close(); err != nil {
			slog.Error("failed to close reader", "error", err)
		}
	}()

	for {
		messages, err := reader.FetchMessages(queues * queueSize)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}

			slog.Error("failed to fetch messages", "error", err)

			continue
		}

		if err := handlerPool.HandleMessages(messages); err != nil {
			slog.Error("failed to handle messages", "error", err)
		}

		if err := reader.CommitMessages(context.Background(), messages...); err != nil {
			slog.Error("failed to commit messages", "error", err)
		}
	}
}
