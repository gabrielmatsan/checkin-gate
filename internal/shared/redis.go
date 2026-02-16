package shared

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Redis struct {
	Client *redis.Client
	logger *zap.Logger
}

func NewRedis(redisURL string, logger *zap.Logger) (*Redis, error) {
	opts, err := redis.ParseURL(redisURL)

	if err != nil {
		logger.Error("failed to parse redis url", zap.Error(err))
		return nil, err
	}

	opts.PoolSize = 10
	opts.MinIdleConns = 5
	opts.MaxRetries = 3
	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 10 * time.Second
	opts.WriteTimeout = 3 * time.Second

	client := redis.NewClient(opts)

	r := &Redis{
		Client: client,
		logger: logger,
	}

	if err := r.HealthCheck(); err != nil {
		logger.Error("failed to connect to redis", zap.Error(err))
		return nil, err
	}

	logger.Info("redis connected successfully")
	return r, nil
}

func (r *Redis) Close() error {
	r.logger.Info("closing redis connection")
	return r.Client.Close()
}

func (r *Redis) HealthCheck() error {
	return r.Client.Ping(context.Background()).Err()
}
