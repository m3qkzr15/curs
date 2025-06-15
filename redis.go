package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456",
		DB:       0,
	})
}

func setRedisKey(redisDB *redis.Client, key, value string) error {
	return redisDB.Set(ctx, key, value, 0).Err()
}

func getRedisKey(redisDB *redis.Client, key string) (string, error) {
	return redisDB.Get(ctx, key).Result()
}

func deleteRedisKey(redisDB *redis.Client, key string) error {
	return redisDB.Del(ctx, key).Err()
}

func getAllRedisKeys(redisDB *redis.Client) (map[string]string, error) {
	var cursor uint64
	data := make(map[string]string)

	for {
		keys, cur, err := redisDB.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			return nil, err
		}
		cursor = cur

		for _, key := range keys {
			value, err := redisDB.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			data[key] = value
		}

		if cursor == 0 {
			break
		}
	}

	return data, nil
}
