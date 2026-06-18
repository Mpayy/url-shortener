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

	return Config
}
