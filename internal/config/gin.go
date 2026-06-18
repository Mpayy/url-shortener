package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewGin(config *viper.Viper, log *logrus.Logger) *gin.Engine {
	app := gin.New()
	app.Use(gin.Recovery())
	return app
}
