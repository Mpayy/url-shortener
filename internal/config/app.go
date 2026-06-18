package config

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type App struct {
	Gin       *gin.Engine
	Logger    *logrus.Logger
	Config    *viper.Viper
	Database  *gorm.DB
	Validator *validator.Validate
	Redis     *RedisClient
}

func NewApp(gin *gin.Engine, logger *logrus.Logger, viper *viper.Viper, db *gorm.DB, validator *validator.Validate, redis *RedisClient) *App {
	app := &App{
		Gin:       gin,
		Logger:    logger,
		Config:    viper,
		Database:  db,
		Validator: validator,
		Redis:     redis,
	}

	return app
}
