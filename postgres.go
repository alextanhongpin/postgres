package postgres

import (
	"database/sql"
	"log"
	"time"

	"github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
)

// New returns a new DB from the config.
func New(cfg *Config, opts ...Options) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.String())
	if err != nil {
		return nil, err
	}

	opt := Option{
		MaxOpenConns:    25,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		Ping:            6,
	}

	for _, modify := range opts {
		modify(&opt)
	}

	if opt.Ping > 0 {
		if err := ping(db, opt.Ping); err != nil {
			return nil, err
		}
	}

	if opt.MigratePath != "" && opt.MigrateTableName != "" {
		if err := makeMigrate(db, opt.MigratePath, opt.MigrateTableName); err != nil {
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

func makeMigrate(db *sql.DB, path, tableName string) error {
	migrate.SetTable(tableName)

	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", path),
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	log.Printf("[migration] Applied %d migrations\n", n)
	return nil
}
