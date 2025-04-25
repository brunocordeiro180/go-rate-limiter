package database

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisDB struct {
	Client *redis.Client
}

func NewRedisConnection(addr, password string) *RedisDB {

	client := redis.NewClient(&redis.Options{Addr: addr, Password: password})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	return &RedisDB{
		Client: client,
	}
}

func (r *RedisDB) Increment(ctx context.Context, key string, expiration int) (int, error) {
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	r.Client.Expire(ctx, key, time.Duration(expiration)*time.Second)
	return int(count), nil
}

func (r *RedisDB) Get(ctx context.Context, key string) (int, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (r *RedisDB) Reset(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
