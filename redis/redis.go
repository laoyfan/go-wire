package redis

import (
	"context"
	"fmt"
	"go-wire/config"
	"go-wire/logger"
	"sync"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Clients map[string]*redis.Client
	log     logger.Logger
}

var ProviderSet = wire.NewSet(NewRedisClients)

func NewRedisClients(cfg *config.Config, log logger.Logger) (*Redis, error) {
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
	return &Redis{Clients: clients, log: log}, nil
}

func (r *Redis) Client(name string) (*redis.Client, error) {
	client, ok := r.Clients[name]
	if !ok {
		return nil, fmt.Errorf("redis 实例 [%s] 不存在", name)
	}
	return client, nil
}

func (r *Redis) Close(ctx context.Context) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for name, client := range r.Clients {
		wg.Add(1)
		go func(name string, client *redis.Client) {
			defer wg.Done()
			if err := client.Close(); err != nil {
				r.log.Error(ctx, fmt.Sprintf("关闭 Redis [%s] 失败", name), logger.Error(err))
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
			} else {
				r.log.Info(ctx, fmt.Sprintf("Redis [%s] 已关闭", name))
			}
		}(name, client)
	}
	wg.Wait()
	return firstErr
}
