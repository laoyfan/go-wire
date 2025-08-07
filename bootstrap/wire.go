//go:build wireinject
// +build wireinject

package bootstrap

import (
	"go-wire/config"
	"go-wire/controller"
	"go-wire/logger"
	"go-wire/redis"
	"go-wire/repo"
	"go-wire/router"
	"go-wire/router/middleware"
	"go-wire/service"

	"github.com/google/wire"
)

func InitApp() (*App, error) {
	wire.Build(
		config.ProviderSet,
		logger.ProviderSet,
		redis.ProviderSet,
		repo.ProviderSet,
		service.ProviderSet,
		controller.ProviderSet,
		middleware.ProviderSet,
		router.ProviderSet,
		ProviderSet,
	)
	return nil, nil
}
