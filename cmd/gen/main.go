package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"kafka-parallel-queues/pkg/kafkakit"

	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaHost := os.Getenv("KAFKA_HOST")
	fmt.Println("in gen/main, kafka ", kafkaHost)
	// 29092
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   []string{kafkaHost},
		Topic:     "numbers",
		BatchSize: 1,
		Balancer:  &kafkakit.SPBalancer{},
	})

	for {
		for _, db := range []string{"A", "B", "C"} {
			for i := 1; i <= 3; i++ {
				err := writer.WriteMessages(context.Background(), kafka.Message{
					Key:   []byte(db),
					Value: []byte(strconv.Itoa(i)),
				})

				if err != nil {
					slog.Error("failed to write messages:", "error", err)
				}
			}
		}
		time.Sleep(5 * time.Minute)
	}

	slog.Info("messages sent")
}
