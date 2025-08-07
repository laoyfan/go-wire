package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewTrans,
	NewApiController,
	wire.Bind(new(RouteRegistrar), new(*ApiController)),
)
