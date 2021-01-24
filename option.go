package postgres

import (
	"time"

	migrate "github.com/rubenv/sql-migrate"
)

// Options for postgres.
type Options struct {
	ConnMaxIdleTime     time.Duration
	ConnMaxLifetime     time.Duration
	MaxIdleConns        int
	MaxOpenConns        int
	MigrationsSource    []migrate.MigrationSource
	MigrationsTableName string
	PingRetries         int
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
func WithMigrationsSource(source migrate.MigrationSource, rest ...migrate.MigrationSource) Option {
	return func(opt *Options) {
		opt.MigrationsSource = append(opt.MigrationsSource, append([]migrate.MigrationSource{source}, rest...)...)
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

// WithConnMaxLifetime overrides the default ConnMaxLifetime of 1
// hour.
func WithConnMaxLifetime(duration time.Duration) Option {
	return func(opt *Options) {
		opt.ConnMaxLifetime = duration
	}
}

// WithConnMaxIdleTime overrides the default ConnMaxIdleTime of 5
// minutes.
func WithConnMaxIdleTime(duration time.Duration) Option {
	return func(opt *Options) {
		opt.ConnMaxIdleTime = duration
	}
}
