package redisClient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"messaging-service/src/types/requests"
	"time"

	"os"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	SetupChannel(channelName string) *redis.PubSub
	PublishToRedisChannel(channelName string, bytes []byte)
	Set(ctx context.Context, key string, value interface{}) error
	SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	GetAPIKey(ctx context.Context, key string) (*requests.APIKey, error)
	GetEmailByPasswordResetToken(ctx context.Context, key string) (string, error)
}

type RedisClient struct {
	Client *redis.Client
}

func connect() (*redis.Client, error) {
	redisURL := os.Getenv("REDIS_URL")
	fmt.Println("CONNECTION STRING IS ", redisURL)
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

func New() *RedisClient {
	fmt.Println("CONNECTING TO REDIS...")
	client, err := connect()
	if err != nil {
		panic(err)
	}
	return &RedisClient{
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

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.Client.Set(ctx, key, p, 0).Result()
	return err
}

func (c *RedisClient) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.Client.Set(ctx, key, p, ttl).Result()
	return err
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	_, err := c.Client.Del(ctx, key).Result()
	return err
}

func (c *RedisClient) GetAPIKey(ctx context.Context, key string) (*requests.APIKey, error) {

	res := c.Client.Get(ctx, key)
	if res.Err() == redis.Nil {
		return nil, nil
	}
	if res.Err() != nil {
		return nil, res.Err()
	}

	apiKey := &requests.APIKey{}
	err := json.Unmarshal([]byte(res.Val()), apiKey)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

func (c *RedisClient) GetEmailByPasswordResetToken(ctx context.Context, key string) (string, error) {
	res := c.Client.Get(ctx, key)
	if res.Err() == redis.Nil {
		return "", nil
	}
	if res.Err() != nil {
		return "", res.Err()
	}

	return string(res.Val()), nil
}
