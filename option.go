package postgres

import (
	"time"

	"github.com/gobuffalo/packr/v2"
)

// Options for postgres.
type Options struct {
	PingRetries         int
	MigrationsSource    *packr.Box
	MigrationsTableName string
	MaxOpenConns        int
	MaxIdleConns        int
	ConnMaxLifetime     time.Duration
}

// Option modify the default Options.
type Option func(*Options)

// WithPing sets the number of times to ping the db every 10 seconds
// before returning error.
func WithPing(n int) Option {
	return func(opt *Options) {
		opt.PingRetries = n
	}
}

// WithMigrationsSource defines the relative path to the folder
// containing migrations, e.g. ./migrations.
func WithMigrationsSource(box *packr.Box) Option {
	return func(opt *Options) {
		opt.MigrationsSource = box
	}
}

// WithMigrationsTableName overrides the rubenv/sql-migrate default
// table called "gorp_migrations".
func WithMigrationsTableName(name string) Option {
	return func(opt *Options) {
		opt.MigrationsTableName = name
	}
}

// WithMaxOpenConns overrides the default MaxOpenConns of 25.
func WithMaxOpenConns(n int) Option {
	return func(opt *Options) {
		opt.MaxOpenConns = n
	}
}

// WithMaxIdleConns overrides the default MaxIdleConns of 25.
func WithMaxIdleConns(n int) Option {
	return func(opt *Options) {
		opt.MaxIdleConns = n
	}
}

// WithConnMaxLifetime overrides the default ConnMaxLifetime of 5
// minutes.
func WithConnMaxLifetime(duration time.Duration) Option {
	return func(opt *Options) {
		opt.ConnMaxLifetime = duration
	}
}
