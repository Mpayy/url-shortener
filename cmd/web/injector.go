//go:build wireinject
// +build wireinject

package main

import (
	"url-shortener/internal/config"
	"url-shortener/internal/delivery/http"
	"url-shortener/internal/delivery/http/middleware"
	"url-shortener/internal/delivery/http/route"
	"url-shortener/internal/repository"
	"url-shortener/internal/usecase"
	"url-shortener/internal/util"

	"github.com/google/wire"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	usecase.NewAuthUsecase,
	http.NewAuthController,
)

var urlSet = wire.NewSet(
	repository.NewUrlRepository,
	usecase.NewUrlUsecase,
	http.NewUrlController,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)

var routeSet = wire.NewSet(
	route.NewRouteConfig,
)

func InitializedApp() *Application {
	wire.Build(
		config.NewViper,
		config.NewLogrus,
		config.NewValidator,
		config.NewGorm,
		config.NewRedis,
		config.NewGin,
		userSet,
		urlSet,
		middlewareSet,
		routeSet,
		util.NewTransaction,
		util.NewTokenUtil,
		config.NewApp,
		NewApplication,
	)
	return nil
}
