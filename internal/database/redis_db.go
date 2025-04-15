package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type RedisDB struct {
	Client  *redis.Client
	Options *redis.Options
}

func NewRedisConnection(addr, password string, db int) *RedisDB {
	return &RedisDB{
		Options: &redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		},
	}
}

func (r *RedisDB) Connect() error {
	r.Client = redis.NewClient(r.Options)
	return nil
}

func (r *RedisDB) Disconnect() error {
	return r.Client.Close()
}
