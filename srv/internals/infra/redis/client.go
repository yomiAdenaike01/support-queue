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

type StreamClient struct {
	Client     *redis.Client
	StreamName string
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

	err = client.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return nil, err
	}

	return &StreamClient{
		Client:     client,
		StreamName: streamName,
	}, nil
}
