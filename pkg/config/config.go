package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    DatabaseURL     string
    AuthServiceAddr string
    AuthzServiceAddr string
    JWTSecret       string
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}