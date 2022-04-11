//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build

package main

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"resource_det_search/internal/conf"
	"resource_det_search/internal/data"
	"resource_det_search/internal/server"
	"resource_det_search/internal/service"
	"resource_det_search/internal/usecase"
)

func initApp(*conf.Data, *zap.SugaredLogger) (*server.Server, func(), error) {
	panic(wire.Build(
		server.ProvideSet, service.ProvideSet, data.ProvideSet, usecase.ProvideSet,
	))
}
