package redis

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go-wire/config"
	"go-wire/logger"
	"sync"
)

type RedisClients struct {
	Clients map[string]*redis.Client
	Logger  logger.Logger
}

var ProviderSet = wire.NewSet(NewRedisClients)

func NewRedisClients(cfg *config.Config, log logger.Logger) (*RedisClients, error) {
	clients := make(map[string]*redis.Client)
	for name, r := range cfg.Redis {
		rdb := redis.NewClient(&redis.Options{
			Addr:     r.Addr,
			Password: r.Password,
			DB:       r.DB,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return nil, fmt.Errorf("redis %s 连接失败: %w", name, err)
		}
		clients[name] = rdb
	}
	return &RedisClients{Clients: clients, Logger: log}, nil
}

func (r *RedisClients) Get(name string) (*redis.Client, error) {
	client, ok := r.Clients[name]
	if !ok {
		return nil, fmt.Errorf("redis 实例 [%s] 不存在", name)
	}
	return client, nil
}

func (r *RedisClients) Close(ctx context.Context) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for name, client := range r.Clients {
		wg.Add(1)
		go func(name string, client *redis.Client) {
			defer wg.Done()
			if err := client.Close(); err != nil {
				r.Logger.Error(ctx, fmt.Sprintf("关闭 Redis [%s] 失败", name), logger.ErrorField(err))
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			} else {
				r.Logger.Info(ctx, fmt.Sprintf("Redis [%s] 已关闭", name))
			}
		}(name, client)
	}
	wg.Wait()
	return firstErr
}
