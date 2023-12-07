package handler

import (
	"context"
	"fmt"
	"log/slog"

	"kafka-parallel-queues/pkg/kafkakit"

	"github.com/go-resty/resty/v2"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

type Pool struct {
	handlerChannels []chan kafka.Message
	errChannel      chan error
	group           *errgroup.Group

	balancer       kafka.Balancer
	channelNumbers []int
}

func NewPool(queues, queueSize int) *Pool {
	var (
		handlerChannels []chan kafka.Message
		channelNumbers  []int
	)

	errChannel := make(chan error, queues*queueSize)
	group, ctx := errgroup.WithContext(context.Background())
	aHandler := New(resty.New())

	for i := 0; i < queues; i++ {
		ch := make(chan kafka.Message, queueSize)
		handlerChannels = append(handlerChannels, ch)
		channelNumbers = append(channelNumbers, i)

		group.Go(func() error {
			for message := range ch {
				errChannel <- aHandler.HandleMessage(ctx, message)
			}

			return nil
		})
	}

	return &Pool{
		handlerChannels: handlerChannels,
		errChannel:      errChannel,
		group:           group,
		channelNumbers:  channelNumbers,
		balancer:        &kafkakit.SHBalancer{},
	}
}

func (handlerPool *Pool) HandleMessages(messages []kafka.Message) error {
	if len(messages) > cap(handlerPool.errChannel) {
		panic(fmt.Sprintf("pool can handle no more %d messages per iteration", cap(handlerPool.errChannel)))
	}

	for _, message := range messages {
		slog.Info("received",
			"number", string(message.Key)+":"+string(message.Value),
		)

		i := handlerPool.balancer.Balance(message, handlerPool.channelNumbers...)

		handlerPool.handlerChannels[i] <- message
	}

	var handlerErr error

	for range messages {
		if err := <-handlerPool.errChannel; err != nil {
			handlerErr = err
		}
	}

	return handlerErr
}

func (handlerPool *Pool) Close() error {
	for _, ch := range handlerPool.handlerChannels {
		close(ch)
	}

	return handlerPool.group.Wait()
}
