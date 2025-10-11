package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type (
	Config struct {
		DataBaseConfig DataBaseConfig `envPrefix:"DB_"`
		HttpConfig     HttpConfig     `envPrefix:"HTTP_"`
	}

	DataBaseConfig struct {
		Host string `env:"URL" envDefault:"postgres://postgres:password@192.168.201.1:5432/UserActivityTracking?sslmode=disable"`
	}

	HttpConfig struct {
		Port int `env:"PORT"  envDefault:"8080"`
	}
)

func GetConfig() (*Config, error) {
	log.Info(fmt.Sprintf("Getting config..."))
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.New("can't parse config")
	}
	return &cfg, nil
}
