package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type App struct {
	Gin      *gin.Engine
	Logger   *logrus.Logger
	Config   *viper.Viper
	Database *gorm.DB
	
}

func NewApp(gin *gin.Engine, logger *logrus.Logger, viper *viper.Viper, db *gorm.DB) *App {
	app := &App{
		Gin:      gin,
		Logger:   logger,
		Config:   viper,
		Database: db,
	}

	return app
}
