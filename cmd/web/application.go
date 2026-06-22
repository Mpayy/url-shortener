package main

import (
	"url-shortener/internal/config"
	"url-shortener/internal/delivery/http/route"
)

type Application struct {
	App         *config.App
	RouteConfig route.RouteConfig
}

func NewApplication(app *config.App, routeConfig route.RouteConfig) *Application {
	return &Application{
		App:         app,
		RouteConfig: routeConfig,
	}
}