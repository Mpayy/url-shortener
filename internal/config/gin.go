package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/toorop/gin-logrus"
)

func NewGin(config *viper.Viper, log *logrus.Logger) *gin.Engine {
	app := gin.New()
	app.Use(ginlogrus.Logger(log), gin.Recovery())
	return app
}
