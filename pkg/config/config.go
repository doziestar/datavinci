package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL         string `mapstructure:"DatabaseURL"`
	AuthServiceAddress  string `mapstructure:"AuthServiceAddress"`
	AuthzServiceAddress string `mapstructure:"AuthzServiceAddress"`
	JWTSecret           string `mapstructure:"JWTSecret"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Enable environment variable reading
	viper.AutomaticEnv()
	viper.BindEnv("DatabaseURL")
	viper.BindEnv("AuthServiceAddress")
	viper.BindEnv("AuthzServiceAddress")
	viper.BindEnv("JWTSecret")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
