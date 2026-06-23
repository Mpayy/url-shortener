package middleware

import (
	"net/http"
	"strings"
	"url-shortener/internal/config"
	"url-shortener/internal/util"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	TokenUtil   util.TokenUtil
	RedisClient config.RedisClient
}

func NewAuthMiddleware(tokenUtil util.TokenUtil, redisClient config.RedisClient) *AuthMiddleware {
	return &AuthMiddleware{TokenUtil: tokenUtil, RedisClient: redisClient}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" || tokenString == "Bearer" {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized")
			return
		}

		exists, err := m.RedisClient.CheckToken(ctx.Request.Context(), tokenString)
		if err != nil || !exists {
			util.ResponseError(ctx, http.StatusUnauthorized, "unauthorized")
			return
		}

		auth, err := m.TokenUtil.ParseToken(tokenString)
		if err != nil {
			util.ResponseError(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		ctx.Set("auth", auth)
		ctx.Next()
	}
}
