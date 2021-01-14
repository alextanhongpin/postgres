package postgres

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host     string `envconfig:"DB_HOST" default:"127.0.0.1"`
	Port     int    `envconfig:"DB_PORT" default:"5432"`
	User     string `envconfig:"DB_USER" default:"john"`
	Password string `envconfig:"DB_PASS" default:"123456"`
	Database string `envconfig:"DB_NAME" default:"development"`
	SSLMode  SSLMode
}

func (c Config) sslmode() string {
	if c.SSLMode == "" {
		return fmt.Sprintf("sslmode=%s", SSLModeDisable)
	}
	return fmt.Sprintf("sslmode=%s", c.SSLMode)
}

func (c Config) String(args ...string) string {
	if len(args) > 0 {
		args = append(args, c.sslmode())
	} else {
		args = []string{c.sslmode()}
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		strings.Join(args, "&"),
	)
}

func NewConfig() (*Config, error) {
	var cfg Config
	var err error
	cfg.Port, err = parseEnvInt("DB_PORT")
	if err != nil {
		return nil, err
	}
	cfg.Host, err = parseEnvString("DB_HOST")
	if err != nil {
		return nil, err
	}
	cfg.User, err = parseEnvString("DB_USER")
	if err != nil {
		return nil, err
	}
	cfg.Password, err = parseEnvString("DB_PASS")
	if err != nil {
		return nil, err
	}
	cfg.Database, err = parseEnvString("DB_NAME")
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func parseEnvInt(env string) (int, error) {
	v := os.Getenv(env)
	if v == "" {
		return 0, fmt.Errorf("%s is required", env)
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func parseEnvString(env string) (string, error) {
	v := os.Getenv(env)
	if v == "" {
		return "", fmt.Errorf("%s is required", env)
	}
	return v, nil
}
