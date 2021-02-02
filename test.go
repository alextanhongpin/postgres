package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/ory/dockertest/v3"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	testContainerImage   = Postgres
	testContainerVersion = "13.1-alpine"

	testHost     = "localhost"
	testPort     = 5432
	testUser     = "root"
	testPassword = "password"
	testDatabase = "test"
)

// NewTestDB returns a new test db with a unique connection string.
func NewTestDB() (*sql.DB, error) {
	return sql.Open("txdb", fmt.Sprintf("tx-%d", time.Now().UnixNano()))
}

type TestOptions struct {
	ConnString       *ConnString
	Container        Container
	MigrationsSource []migrate.MigrationSource
}

type TestOption func(*TestOptions)

// WithContainer overrides the current container.
func WithContainer(ctn Container) TestOption {
	return func(opt *TestOptions) {
		opt.Container = ctn
	}
}

// WithConnString sets the test connection string.
func WithConnString(cs *ConnString) TestOption {
	return func(opt *TestOptions) {
		opt.ConnString = cs
	}
}

// WithTestMigrationsSource sets the path to the migrations folder.
func WithTestMigrationsSource(src migrate.MigrationSource, rest ...migrate.MigrationSource) TestOption {
	return func(opt *TestOptions) {
		opt.MigrationsSource = append(opt.MigrationsSource, append([]migrate.MigrationSource{src}, rest...)...)
	}
}

// InitTestDB setups a dockertest postgres container.
func InitTestDB(opts ...TestOption) (*sql.DB, func() error) {
	opt := TestOptions{
		Container: Container{
			Image:   testContainerImage,
			Version: testContainerVersion,
		},
		ConnString: &ConnString{
			Host:     testHost,
			Port:     testPort,
			User:     testUser,
			Password: testPassword,
			Database: testDatabase,
		},
	}

	for _, modify := range opts {
		modify(&opt)
	}

	var db *sql.DB
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("connect docker failed: %s", err)
	}

	cs := opt.ConnString
	cs.Host = testHost

	// Pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: opt.Container.Image,
		Tag:        opt.Container.Version,
		Cmd:        []string{"-c", "fsync=off"},
		Env: []string{
			fmt.Sprintf("POSTGRES_DB=%s", cs.Database),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cs.Password),
			fmt.Sprintf("POSTGRES_USER=%s", cs.User),
		},
	})
	if err != nil {
		log.Fatalf("start docker failed: %s", err)
	}

	// Hard kill the container in 60 seconds.
	_ = resource.Expire(60)

	// Exponential backoff-retry, because the application in the container might
	// not be ready to accept connections yet
	if err := pool.Retry(func() error {
		port, err := strconv.Atoi(resource.GetPort("5432/tcp"))
		if err != nil {
			return err
		}

		// Assign dynamic port.
		cs.Port = port
		db, err = sql.Open(Postgres, cs.String())
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("connect failed: %s", err)
	}

	// Perform migration.
	for _, src := range opt.MigrationsSource {
		if err := Migrate(db, src); err != nil {
			log.Fatalf("migration failed: %s", err)
		}
	}

	txdb.Register("txdb", Postgres, cs.String())

	return db, func() error {
		if err := db.Close(); err != nil {
			return err
		}

		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resource); err != nil {
			return err
		}

		return nil
	}
}
