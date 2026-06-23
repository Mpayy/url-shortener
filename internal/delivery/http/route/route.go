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
	// Public routes
	public := r.App.Group("/api/v1")
	public.POST("/auth/register", r.AuthController.Register)
	public.POST("/auth/login", r.AuthController.Login)

	// Protected routes
	protected := r.App.Group("/api/v1", r.AuthMiddleware.Authenticate())
	protected.DELETE("/auth/logout", r.AuthController.Logout)
	// protected.POST("/urls", r.UrlController.Create)
	// protected.GET("/urls", r.UrlController.GetAll)
	// protected.DELETE("/urls/:short_code", r.UrlController.Delete)

	// Public redirect
	// r.App.GET("/:short_code", r.UrlController.Redirect)
}
