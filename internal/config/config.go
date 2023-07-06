package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Env  string
	Port string

	RedisURI string `mapstructure:"REDIS_URI"`
	MongoURI string `mapstructure:"MONGO_URI"`

	JwtSecret string `mapstructure:"JWT_SECRET"`
	JwtIssuer string `mapstructure:"JWT_ISSUER"`
	JwtExp    int    `mapstructure:"JWT_EXP"`
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
