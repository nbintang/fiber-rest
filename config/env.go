package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Env struct {
	DatabaseURL      string `mapstructure:"DATABASE_URL"`
	DatabaseHost     string `mapstructure:"DATABASE_HOST"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName     string `mapstructure:"DATABASE_NAME"`
	DatabasePort     int    `mapstructure:"DATABASE_PORT"`
	AppEnv           string `mapstructure:"APP_ENV"`
	AppAddr          string `mapstructure:"APP_ADDR"`
}

func NewEnv() (Env, error) {
	viper.SetDefault("APP_ENV", "development")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	var env Env

	if err := viper.Unmarshal(&env); err != nil {
		return Env{}, err
	}
	if env.DatabaseURL == "" {
		return Env{}, fmt.Errorf("DATABASE_URL is required")
	}

	return env, nil
}
