package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogrus(config *viper.Viper) *logrus.Logger {
	log := logrus.New()

	// jika mau lognya di tulis dalam bentuk file
	// file, _ := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// log.SetOutput(file)

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	if levelStr := config.GetString("LOG_LEVEL"); levelStr != "" {
		if level, err := logrus.ParseLevel(levelStr); err == nil {
			log.SetLevel(level)
		}
	}

	return log
}
