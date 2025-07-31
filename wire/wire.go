//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"go-wire/bootstrap"
	"go-wire/config"
	"go-wire/logger"
	"go-wire/redis"
	"go-wire/router"
)

func InitApp() (*bootstrap.App, error) {
	wire.Build(
		config.ProviderSet,
		logger.ProviderSet,
		redis.ProviderSet,
		router.ProviderSet,
		bootstrap.ProviderSet,
	)
	return nil, nil
}
