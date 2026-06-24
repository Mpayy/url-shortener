package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	Config := viper.New()

	Config.SetConfigName("config")
	Config.SetConfigType("env")
	Config.AddConfigPath("./")
	Config.AddConfigPath("../../")

	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	secret := Config.GetString("JWT_SECRET_KEY")
	if len(secret) < 32 {
		panic("JWT_SECRET_KEY must be at least 32 characters")
	}

	return Config
}
