package handler

import (
	"context"
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	httpClient *resty.Client
}

func New(httpClient *resty.Client) *Handler {
	return &Handler{
		httpClient: httpClient,
	}
}

func (h *Handler) HandleMessage(ctx context.Context, message kafka.Message) error {
	dbURL := getDBUrl(string(message.Key))

	_, err := h.httpClient.R().
		Post(dbURL + string(message.Key) + ":" + string(message.Value))

	if err != nil {
		slog.Error("failed to write a number in db",
			"url", dbURL,
			"error", err,
		)
	}

	return nil
}

func getDBUrl(messageKey string) string {
	switch messageKey {
	case "A":
		return "http://db_a/api/write/"
	case "B":
		return "http://db_b/api/write/"
	case "C":
		return "http://db_c/api/write/"
	}

	return ""
}
