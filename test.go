package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/gobuffalo/packr/v2"
	"github.com/ory/dockertest/v3"
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
	Container        Container
	MigrationsSource *packr.Box
	ConnString       *ConnString
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
func WithTestMigrationsSource(box *packr.Box) TestOption {
	return func(opt *TestOptions) {
		opt.MigrationsSource = box
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
	resource, err := pool.Run(
		opt.Container.Image,
		opt.Container.Version,
		[]string{
			fmt.Sprintf("POSTGRES_DB=%s", cs.Database),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cs.Password),
			fmt.Sprintf("POSTGRES_USER=%s", cs.User),
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
	if opt.MigrationsSource != nil {
		if err := makeMigrate(db, opt.MigrationsSource); err != nil {
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
