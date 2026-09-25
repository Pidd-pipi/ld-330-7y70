package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	DBHost      string `env:"DB_HOST" envDefault:"localhost"`
	DBPort      string `env:"DB_PORT" envDefault:"5432"`
	DBName      string `env:"DB_NAME" envDefault:"gbemr"`
	DBUser      string `env:"DB_USER" envDefault:"gbemr"`
	DBPassword  string `env:"DB_PASSWORD" envDefault:"gbemr_password"`
	DBSSLMode   string `env:"DB_SSLMODE" envDefault:"disable"`
	DatabaseURL string `env:"DATABASE_URL"`
	JWTSecret   string `env:"JWT_SECRET" envDefault:"change-this-development-secret"`
}

func Load() (Config, error) { return env.ParseAs[Config]() }
func (c Config) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}
