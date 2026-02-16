package config

import "github.com/caarlos0/env/v11"

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

type Config struct {
	DatabaseURL        string      `env:"DATABASE_URL"`
	Port               int         `env:"PORT" envDefault:"8080"`
	JWTSecret          string      `env:"JWT_SECRET,required"`
	Env                Environment `env:"ENV" envDefault:"development"`
	GoogleClientID     string      `env:"GOOGLE_CLIENT_ID,required"`
	GoogleClientSecret string      `env:"GOOGLE_CLIENT_SECRET,required"`
	GoogleRedirectURL  string      `env:"GOOGLE_REDIRECT_URL" envDefault:"http://localhost:8080/auth/google/callback"`
	RedisURL   string `env:"REDIS_URL,required"`
	ResendKey  string `env:"RESEND_KEY,required"`
	ResendFrom string `env:"RESEND_FROM" envDefault:"gabriel@laboratorio-de-pesquisa-de-engenharia-de-software.com"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
