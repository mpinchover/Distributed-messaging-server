package redisClient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func New() RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return RedisClient{
		Client: client,
	}
	// There is no error because go-redis automatically reconnects on error.
}

func (c *RedisClient) Save(channelName string) <-chan *redis.Message {

	// subscribe to the channel
	subscriber := c.Client.Subscribe(context.Background(), channelName)
	ch := subscriber.Channel()
	return ch
}

// // pass in the identifier of the channel so Redis can perform pub/sub
func (c *RedisClient) SetupChannel(channelName string) *redis.PubSub {
	subscriber := c.Client.Subscribe(context.Background(), channelName)
	return subscriber
}

// publish message to redis channel
func (c *RedisClient) PublishToRedisChannel(channelName string, bytes []byte) {
	c.Client.Publish(context.Background(), channelName, bytes)
}
