package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) (*RedisCache, error) {
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{client}, nil
}

func (c *RedisCache) Get(key string) (string, error) {
	val, err := c.client.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (c *RedisCache) Set(key string, value string, expiration time.Duration) error {
	return c.client.Set(key, value, expiration).Err()
}

func (c *RedisCache) Delete(key string){

}