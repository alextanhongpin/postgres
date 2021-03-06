package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	migrate "github.com/rubenv/sql-migrate"
)

const Postgres = "postgres"

const (
	// Migrations.
	migrationsTableName = "migrations"

	// Ping.
	pingRetries = 6
	pingTimeout = 10 * time.Second

	// Connectivity.
	maxOpenConns    = 25
	maxIdleConns    = 25
	connMaxLifetime = 1 * time.Hour
	connMaxIdleTime = 5 * time.Minute
)

func New(connString string, options ...Option) (*sql.DB, error) {
	db, err := sql.Open(Postgres, connString)
	if err != nil {
		return nil, err
	}

	opts := Options{
		MaxOpenConns:        maxOpenConns,
		MaxIdleConns:        maxIdleConns,
		ConnMaxIdleTime:     connMaxIdleTime,
		ConnMaxLifetime:     connMaxLifetime,
		PingRetries:         pingRetries,
		MigrationsTableName: migrationsTableName,
	}

	for _, opt := range options {
		opt(&opts)
	}

	if opts.PingRetries > 0 {
		if err := ping(db, opts.PingRetries); err != nil {
			return nil, err
		}
	}

	migrate.SetTable(opts.MigrationsTableName)
	for _, src := range opts.MigrationsSource {
		if err := Migrate(db, src); err != nil {
			return nil, err
		}
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(opts.MaxOpenConns)
	db.SetMaxIdleConns(opts.MaxIdleConns)
	db.SetConnMaxLifetime(opts.ConnMaxLifetime)
	db.SetConnMaxIdleTime(opts.ConnMaxIdleTime)

	return db, nil
}

func ping(db *sql.DB, retry int) error {
	var err error
	for i := 0; i < retry; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(pingTimeout)
	}
	return err
}

func Migrate(db *sql.DB, migrations migrate.MigrationSource) error {
	n, err := migrate.Exec(db, Postgres, migrations, migrate.Up)
	if err != nil {
		return err
	}
	log.Printf("migrations applied: %d\n", n)
	return nil
}
