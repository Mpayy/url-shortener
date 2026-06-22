package route

import (
	"url-shortener/internal/delivery/http"
	"url-shortener/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App            *gin.Engine
	AuthMiddleware *middleware.AuthMiddleware
	AuthController http.AuthController
}

func NewRouteConfig(app *gin.Engine, authMiddleware *middleware.AuthMiddleware, authController http.AuthController) RouteConfig {
	return RouteConfig{App: app, AuthMiddleware: authMiddleware, AuthController: authController}
}

func (r *RouteConfig) Setup() {
	auth := r.App.Group("/api/v1/auth")
	auth.POST("/register", r.AuthController.Register)
}
