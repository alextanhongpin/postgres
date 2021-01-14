package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/ory/dockertest/v3"
)

// NewTestDB returns a new test db with a unique connection string.
func NewTestDB() (*sql.DB, error) {
	return sql.Open("txdb", fmt.Sprintf("tx-%d", time.Now().UnixNano()))
}

type TestOption struct {
	DockerImage      DockerImage
	MigrationsSource string
	ConnString       *ConnString
}

type TestOptionModifier func(*TestOption)

func WithDockerImage(di DockerImage) TestOptionModifier {
	return func(opt *TestOption) {
		opt.DockerImage = di
	}
}

func WithConnString(cs *ConnString) TestOptionModifier {
	return func(opt *TestOption) {
		opt.ConnString = cs
	}
}

func WithTestMigrationsSource(src string) TestOptionModifier {
	return func(opt *TestOption) {
		opt.MigrationsSource = src
	}
}

// SetupTestDB setups a dockertest postgres container.
func SetupTestDB(opts ...TestOptionModifier) (*sql.DB, func()) {
	opt := TestOption{
		DockerImage: DockerImage{
			Container: "postgres",
			Version:   "13.1-alpine",
		},
		ConnString: &ConnString{
			Host:     "localhost",
			Port:     5432,
			User:     "root",
			Password: "password",
			Database: "test",
		},
	}

	for _, modify := range opts {
		modify(&opt)
	}

	var db *sql.DB
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	cs := opt.ConnString
	cs.Host = "localhost"

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run(
		opt.DockerImage.Container,
		opt.DockerImage.Version,
		[]string{
			fmt.Sprintf("POSTGRES_DB=%s", cs.Database),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", cs.Password),
			fmt.Sprintf("POSTGRES_USER=%s", cs.User),
		})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Hard kill the container in 60 seconds.
	_ = resource.Expire(60)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		port, err := strconv.Atoi(resource.GetPort("5432/tcp"))
		if err != nil {
			return err
		}
		// Port changes dynamically.
		cs.Port = port
		db, err = sql.Open("postgres", cs.String())
		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Perform migration.
	if opt.MigrationsSource != "" {
		if err := makeMigrate(db, opt.MigrationsSource); err != nil {
			log.Fatalf("test migration failed: %v", err)
		}
	}

	txdb.Register("txdb", "postgres", cs.String())

	return db, func() {
		if err := db.Close(); err != nil {
			log.Printf("Could not close db: %s", err)
		}

		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
