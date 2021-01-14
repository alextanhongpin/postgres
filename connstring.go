package postgres

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ConnString struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Params   map[string]string
}

func (c ConnString) String() string {
	var params []string
	if _, ok := c.Params["sslmode"]; !ok {
		params = append(params, "sslmode=disable")
	}
	for k, v := range c.Params {
		params = append(params, fmt.Sprintf("%v=%v", k, v))
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		strings.Join(params, "&"),
	)
}

func NewConnString() (*ConnString, error) {
	var cfg ConnString
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
