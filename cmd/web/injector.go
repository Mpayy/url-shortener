//go:build wireinject
// +build wireinject

package main

import (
	"url-shortener/internal/config"

	"github.com/google/wire"
)

func InitializedApp() *config.App {
	wire.Build(
		config.NewViper,
		config.NewLogrus,
		config.NewValidator,
		config.NewGorm,
		config.NewRedis,
		config.NewGin
		config.NewApp,
	)
	return nil
}
