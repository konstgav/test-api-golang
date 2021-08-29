package redis

import (
	"context"
	"encoding/json"
	"log"
	"test-api-golang/model"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	Client     *redis.Client
	expiration time.Duration
}

func NewCache(addr string, passw string, db int, expiration time.Duration) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passw,
		DB:       db,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Println("Connect to redis, ", pong)
	return &Cache{Client: rdb}
}

var ctx = context.Background()

func (c *Cache) Set(key string, value interface{}) error {
	log.Println("Set to redis, id", key)
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Client.Set(ctx, key, json, c.expiration).Err()
}

func (c *Cache) Get(key string) (interface{}, error) {
	log.Println("Get from redis, id", key)
	value, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var product model.Product
	err = json.Unmarshal([]byte(value), &product)
	if err != nil {
		panic(err)
	}
	return &product, nil
}
