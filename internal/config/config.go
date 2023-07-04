package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env  string
	Port string

	RedisDNS string `mapstructure:"REDIS_DNS"`
}

func Load(config *AppConfig) error {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&config)
	return err
}
