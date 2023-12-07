package kafkakit

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	*kafka.Reader

	fetchMessageChannel chan messageOrError
}

func NewReader(fetchContext context.Context, cfg kafka.ReaderConfig) *Reader {
	kafkaReader := kafka.NewReader(cfg)

	return &Reader{
		Reader:              kafkaReader,
		fetchMessageChannel: newFetchMessageChannel(fetchContext, kafkaReader),
	}
}

func (r *Reader) FetchMessages(amount int) ([]kafka.Message, error) {
	const waitTimeout = 100 * time.Millisecond

	messages := make([]kafka.Message, 0, amount)

	ticker := time.NewTicker(waitTimeout)

	for {
		select {
		case <-ticker.C:
			if len(messages) > 1 {
				return messages, nil
			}

			ticker.Reset(waitTimeout)
		case msgOrErr := <-r.fetchMessageChannel:
			if msgOrErr.err != nil {
				return nil, msgOrErr.err
			}

			messages = append(messages, msgOrErr.message)

			if len(messages) == amount {
				return messages, nil
			}
		}
	}
}

type messageOrError struct {
	message kafka.Message
	err     error
}

func newFetchMessageChannel(ctx context.Context, kafkaReader *kafka.Reader) chan messageOrError {
	fetchMessageChan := make(chan messageOrError)

	go func() {
		for {
			message, err := kafkaReader.FetchMessage(ctx)
			if err != nil {
				fetchMessageChan <- messageOrError{err: err}

				if errors.Is(err, context.Canceled) {
					close(fetchMessageChan)

					return
				}

				continue
			}

			fetchMessageChan <- messageOrError{message: message}
		}
	}()

	return fetchMessageChan
}
