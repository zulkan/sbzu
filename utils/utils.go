package utils

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

func ReadJSON[T any](s string) *T {
	var res T
	err := json.Unmarshal([]byte(s), &res)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &res
}

func ToJSON(data any) string {
	dataStr, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(dataStr)
}

func ToInt64(input string) int64 {
	val, _ := strconv.ParseInt(input, 10, 64)
	return val
}

func WriteRedis(redisClient *redis.Client, key, value string, duration time.Duration) error {
	return redisClient.Set(key, value, duration).Err()
}

func ReadRedis(redisClient *redis.Client, key string) (string, error) {
	val, err := redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func RemoveRedis(redisClient *redis.Client, key string) error {
	_, err := redisClient.Del(key).Result()
	return err
}
