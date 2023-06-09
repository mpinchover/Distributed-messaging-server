package redisClient

import (
	"context"
	"errors"

	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func connect() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	var statusCode string
	status := client.Ping(ctx)
	statusCode = status.Val()
	if statusCode != "PONG" {
		return nil, errors.New("could not connect to redis")
	}
	return client, nil
}

func New() RedisClient {
	client, err := connect()
	if err != nil {
		panic(err)
	}
	return RedisClient{
		Client: client,
	}
}

// pass in the identifier of the channel so Redis can perform pub/sub
func (c *RedisClient) SetupChannel(channelName string) *redis.PubSub {
	subscriber := c.Client.Subscribe(context.Background(), channelName)
	return subscriber
}

// publish message to redis channel
func (c *RedisClient) PublishToRedisChannel(channelName string, bytes []byte) {
	c.Client.Publish(context.Background(), channelName, bytes)
}
