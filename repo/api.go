package repo

import (
	"go-wire/logger"
	"go-wire/redis"

	"github.com/gin-gonic/gin"
)

type ApiRepo struct {
	redis *redis.Redis
	log   logger.Logger
}

func NewApiRepo(redis *redis.Redis, log logger.Logger) *ApiRepo {
	return &ApiRepo{redis: redis, log: log}
}

func (r *ApiRepo) Test(ctx *gin.Context, id string) (string, error) {
	redisClient, err := r.redis.Client("default")
	if err != nil {
		r.log.Error(ctx, "redis nil", logger.Error(err))
		return "", err
	}
	userValue, err := redisClient.Get(ctx, id).Result()
	if err != nil {
		r.log.Error(ctx, "未获取到用户", logger.Error(err))
		return "", err
	}
	return userValue, nil
}
