package helpers

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisCacherAdapter(c *redis.Client) *RedisCacherAdapter {
	return &RedisCacherAdapter{c: c}
}

type RedisCacherAdapter struct {
	c *redis.Client
}

func RedisConnect() (*RedisCacherAdapter, error) {
	redisDBSelect, err := strconv.Atoi(os.Getenv("REDIS_DB_SELECT"))
	if err != nil {
		log.Println("main() redisDBSelect is not an integer, failed to parse: " + err.Error())
		return nil, err
	}
	rc := redis.NewClient(&redis.Options{
		Addr:       os.Getenv("REDIS_URL"),
		Password:   os.Getenv("REDIS_PASSWORD"),
		DB:         redisDBSelect,
		MaxRetries: 3,
	})
	return NewRedisCacherAdapter(rc), nil
}

func (r *RedisCacherAdapter) RedisSet(ctx context.Context, key string, val string, dur time.Duration) (err error) {
	return r.c.Set(ctx, key, val, dur).Err()
}
func (r *RedisCacherAdapter) RedisGet(ctx context.Context, key string) (val string, err error) {
	v, err := r.c.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return "", err
		}
	}

	return v, nil
}
