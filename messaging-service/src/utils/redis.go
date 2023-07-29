package utils

import (
	"context"
	"encoding/json"
	"log"
	redisClient "messaging-service/src/redis"
	"messaging-service/src/types/requests"

	"github.com/redis/go-redis/v9"
)

func SetClientConnectionToRedis(ctx context.Context, client *redisClient.RedisClient, connection *requests.Connection) error {
	return client.Set(ctx, connection.UserUUID, connection)
}

// subscribe to the channel
func SubscribeToChannel(subscriber *redis.PubSub, fn func(event string) error) {
	for redisMsg := range subscriber.Channel() {
		err := fn(redisMsg.Payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// pass in the identifier of the channel so Redis can perform pub/sub
// TODO – use handler context
func SetupChannel(c *redisClient.RedisClient, channelName string) *redis.PubSub {
	subscriber := c.Client.Subscribe(context.Background(), channelName)
	return subscriber
}

// TODO use context from handler
func PublishToRedisChannel(c *redisClient.RedisClient, channelName string, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	res := c.Client.Publish(context.Background(), channelName, bytes)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
