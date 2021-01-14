package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
)

func New(connString string, opts ...OptionModifier) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	opt := Option{
		MaxOpenConns:        25,
		MaxIdleConns:        25,
		ConnMaxLifetime:     5 * time.Minute,
		Ping:                6,
		MigrationsTableName: "migrations",
	}

	for _, modify := range opts {
		modify(&opt)
	}

	if opt.Ping > 0 {
		if err := ping(db, opt.Ping); err != nil {
			return nil, err
		}
	}

	migrate.SetTable(opt.MigrationsTableName)
	if opt.MigrationsSource != "" {
		if err := makeMigrate(db, opt.MigrationsSource); err != nil {
			return nil, err
		}
	}

	// https://www.alexedwards.net/blog/configuring-sqldb
	db.SetMaxOpenConns(opt.MaxOpenConns)
	db.SetMaxIdleConns(opt.MaxIdleConns)
	db.SetConnMaxLifetime(opt.ConnMaxLifetime)

	return db, nil
}

func ping(db *sql.DB, retry int) error {
	var err error
	for i := 0; i < retry; i++ {
		err = db.Ping()
		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	return err
}

func makeMigrate(db *sql.DB, src string) error {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", src),
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	log.Printf("[migration] Applied %d migrations\n", n)
	return nil
}
