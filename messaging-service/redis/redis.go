package redisClient

import (
	"context"

	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func New() RedisClient {
	// log.Println(os.Getenv("REDIS_URL"))
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
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
