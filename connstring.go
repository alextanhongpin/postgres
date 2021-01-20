package postgres

import (
	"fmt"
	"strings"

	"github.com/alextanhongpin/postgres/internal/env"
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

// NewConnString creates a connection string from environment variables.
func NewConnString() *ConnString {
	return &ConnString{
		Host:     env.MustString("DB_HOST"),
		Port:     env.MustInt("DB_PORT"),
		User:     env.MustString("DB_USER"),
		Password: env.MustString("DB_PASS"),
		Database: env.MustString("DB_NAME"),
	}
}
