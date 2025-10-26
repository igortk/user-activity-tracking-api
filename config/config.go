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
		CronConfig     CronConfig     `envPrefix:"CRON_"`
	}

	DataBaseConfig struct {
		Host            string `env:"URL" envDefault:"postgres://postgres:password@192.168.201.1:5432/UserActivityTracking?sslmode=disable"`
		MaxOpenConns    int64  `env:"MAX_OPEN_CONN" envDefault:"50"`
		MaxIdleConns    int64  `env:"MAX_IDLE_CONN" envDefault:"15"`
		ConnMaxLifetime int64  `env:"CONN_MAX_LIFETIME_MINUTE" envDefault:"15"`
	}

	HttpConfig struct {
		Port              int64      `env:"PORT"  envDefault:"8080"`
		ReadTimeout       int64      `env:"READ_TIMEOUT"  envDefault:"5"`
		ReadHeaderTimeout int64      `env:"READ_HEADER_TIMEOUT"  envDefault:"2"`
		WriteTimeout      int64      `env:"WRITE_TIMEOUT"  envDefault:"10"`
		IdleTimeout       int64      `env:"IDLE_TIMEOUT"  envDefault:"120"`
		CorsConfig        CorsConfig `envPrefix:"CORS_"`
	}

	CronConfig struct {
		Tab CronTab `envPrefix:"TAB_"`
	}

	CronTab struct {
		TabCountUsersEventTask string `env:"TAB_COUNT_USERS_EVENT_TASK" envDefault:"* * * * *"`
	}

	CorsConfig struct {
		AllowedOrigins   string `env:"ALLOWED_ORIGINS" envDefault:"*"`
		AllowMethods     string `env:"ALLOWED_METHODS" envDefault:"GET,POST"`
		AllowHeaders     string `env:"ALLOWED_HEADERS" envDefault:"Content-Type"`
		MaxAgeHoursCache int64  `env:"MAX_AGE_HOURS_CACHE" envDefault:"12"`
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
