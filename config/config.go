package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	LimitByIPMaxRPS         int64  `mapstructure:"LIMIT_BY_IP_MAX_RPS"`
	LimitByIPBlockTimeMs    int64  `mapstructure:"LIMIT_BY_IP_BLOCK_TIME_MS"`
	LimitByTokenMaxRPS      int64  `mapstructure:"LIMIT_BY_TOKEN_MAX_RPS"`
	LimitByTokenBlockTimeMs int64  `mapstructure:"LIMIT_BY_TOKEN_BLOCK_TIME_MS"`
	WebServerPort           string `mapstructure:"WEB_SERVER_PORT"`
	RedisAddr               string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	viper.AutomaticEnv()
	cfg.LimitByIPMaxRPS = viper.GetInt64("LIMIT_BY_IP_MAX_RPS")
	cfg.LimitByIPBlockTimeMs = viper.GetInt64("LIMIT_BY_IP_BLOCK_TIME_MS")
	cfg.LimitByTokenMaxRPS = viper.GetInt64("LIMIT_BY_TOKEN_MAX_RPS")
	cfg.LimitByTokenBlockTimeMs = viper.GetInt64("LIMIT_BY_TOKEN_BLOCK_TIME_MS")
	cfg.WebServerPort = viper.GetString("WEB_SERVER_PORT")
	cfg.RedisAddr = viper.GetString("REDIS_ADDRESS")
	if cfg.WebServerPort == "" || cfg.RedisAddr == "" {
		return nil, os.ErrInvalid
	}
	return &cfg, nil
}
