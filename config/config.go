package config

import (
	"github.com/spf13/viper"
)

type config struct {
	RateLimitIP    int    `mapstructure:"RATE_LIMIT_IP"`
	RateLimitToken int    `mapstructure:"RATE_LIMIT_TOKEN"`
	RateDuration   int    `mapstructure:"RATE_DURATION"`
	RedisAddr      string `mapstructure:"REDIS_ADDR"`
	RedisPassword  string `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig(path string) (*config, error) {
	var cfg *config

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, err
}
