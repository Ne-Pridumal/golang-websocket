package rs

import (
	"context"
	"golang-websocket-chat/internal/config"

	"github.com/redis/go-redis/v9"
)

func New(opt config.Redis) (*redis.Client, error) {
	ctx := context.Background()
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     opt.Address,
			Password: opt.Password,
			DB:       0,
		})
	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}

	return rdb, nil
}
