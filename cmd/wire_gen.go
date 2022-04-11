// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"go.uber.org/zap"
	"resource_det_search/internal/conf"
	"resource_det_search/internal/data"
	"resource_det_search/internal/server"
	"resource_det_search/internal/service"
	"resource_det_search/internal/usecase"
)

// Injectors from wire.go:

func initApp(confData *conf.Data, sugaredLogger *zap.SugaredLogger) (*server.Server, func(), error) {
	dataData, cleanup, err := data.NewData(confData, sugaredLogger)
	if err != nil {
		return nil, nil, err
	}
	iUserRepo := data.NewUserRepo(dataData)
	iUserUsecase := usecase.NewUserUsecase(iUserRepo)
	userService := service.NewUserService(iUserUsecase, sugaredLogger)
	iDimensionRepo := data.NewDimensionRepo(dataData)
	iDimensionUsecase := usecase.NewDimensionUsecase(iDimensionRepo)
	dimensionService := service.NewDimensionService(iDimensionUsecase, sugaredLogger)
	serverServer := server.NewServer(userService, dimensionService)
	return serverServer, func() {
		cleanup()
	}, nil
}