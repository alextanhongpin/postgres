package postgres

import "time"

// Option for postgres.
type Option struct {
	Ping                int
	MigrationsSource    string
	MigrationsTableName string
	MaxOpenConns        int
	MaxIdleConns        int
	ConnMaxLifetime     time.Duration
}

// OptionModifier modify the default Option.
type OptionModifier func(*Option)

// WithPing sets the number of times to ping the db every 10 seconds
// before returning error.
func WithPing(n int) OptionModifier {
	return func(opt *Option) {
		opt.Ping = n
	}
}

// WithMigrationsSource defines the relative path to the folder
// containing migrations, e.g. ./migrations.
func WithMigrationsSource(src string) OptionModifier {
	return func(opt *Option) {
		opt.MigrationsSource = src
	}
}

// WithMigrationsTableName overrides the rubenv/sql-migrate default
// table called "gorp_migrations".
func WithMigrationsTableName(name string) OptionModifier {
	return func(opt *Option) {
		opt.MigrationsTableName = name
	}
}

// WithMaxOpenConns overrides the default MaxOpenConns of 25.
func WithMaxOpenConns(n int) OptionModifier {
	return func(opt *Option) {
		opt.MaxOpenConns = n
	}
}

// WithMaxIdleConns overrides the default MaxIdleConns of 25.
func WithMaxIdleConns(n int) OptionModifier {
	return func(opt *Option) {
		opt.MaxIdleConns = n
	}
}

// WithConnMaxLifetime overrides the default ConnMaxLifetime of 5
// minutes.
func WithConnMaxLifetime(duration time.Duration) OptionModifier {
	return func(opt *Option) {
		opt.ConnMaxLifetime = duration
	}
}
