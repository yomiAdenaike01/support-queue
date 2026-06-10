package redisinfra

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yomiAdenaike01/support-queue/internals/config"
)

type Streamname = string

const (
	RESOLVED_TICKET_STREAM Streamname = "TICKETS:RESOLVED_STREAM"
	TICKET_CREATED         Streamname = "TICKETS_STREAM"
)

type StreamClient struct {
	client     *redis.Client
	streamName Streamname
}

type PushEventType string

const (
	PUSHEVENT_TICKET_SUBMITTED = "TICKET_SUBMITTED"
)

type PushEvent struct {
	StreamName string                 `json:"stream_name"`
	Values     map[string]interface{} `json:"values"`
	EventType  PushEventType          `json:"type"`
}

func (s *StreamClient) Push(ctx context.Context, event PushEvent) error {
	_, err := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: event.StreamName,
		Values: event.Values,
	}).Result()
	return err
}

func NewStreamClient(ctx context.Context, config *config.Config) (*StreamClient, error) {
	streamName := config.GetEnvOrFail("STREAM_NAME")
	consumerGroup := config.GetEnvOrFail("STREAM_CONSUMER_GROUP")

	u, err := url.Parse(config.GetEnvOrFail("REDIS_ADDR"))
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	pingCtx, cancel := context.WithTimeout(ctx, time.Duration(time.Second*30))
	defer cancel()

	res, err := client.Ping(pingCtx).Bytes()
	log.Println(string(res))
	if err != nil {
		return nil, err
	}

	for _, streamName := range []Streamname{TICKET_CREATED, RESOLVED_TICKET_STREAM} {
		err = client.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err()
		if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
			return nil, err
		}
	}

	return &StreamClient{
		client:     client,
		streamName: streamName,
	}, nil
}
