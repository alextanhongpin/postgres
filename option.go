package postgres

import "time"

type Option struct {
	Ping             int
	MigratePath      string
	MigrateTableName string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
}

type Options func(*Option)

func WithPing(n int) Options {
	return func(opt *Option) {
		opt.Ping = n
	}
}

func WithMigrate(path, tableName string) Options {
	return func(opt *Option) {
		opt.MigratePath = path
		opt.MigrateTableName = tableName
	}
}

func WithMaxOpenConns(n int) Options {
	return func(opt *Option) {
		opt.MaxOpenConns = n
	}
}

func WithConnMaxLifetime(duration time.Duration) Options {
	return func(opt *Option) {
		opt.ConnMaxLifetime = duration
	}
}

func WithMaxIdleConns(n int) Options {
	return func(opt *Option) {
		opt.MaxIdleConns = n
	}
}
