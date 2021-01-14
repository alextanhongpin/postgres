# postgres

A simple wrapper around basic operations with Postgres in Golang.


## Usage

Creating a new connection string:
```go
cs := ConnString{
	Host: "127.0.0.1",
	Port: 5342,
	User: "john",
	Password: "123456",
	Database: "test",
	Params: map[string]string{
		"sslmode": "disable",  // Automatically added.
	},
}
connStr := cs.String()
```

# Reading from environment variables

This reads the environment variables `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER` and `DB_PASS`. It will return an error if any of that is missing.

```go
cs, err := NewConnString()
if err != nil {
	log.Fatal(err)
}
connStr := cs.String()
```

# New DB Instance

Creates a new database instance:
```go
db, err := New(cs.String())
if err != nil {
	log.Fatal(err)
}
defer db.Close()
```

It also takes additional params to define the migrations source using packr2, so that you can run the migration when the application starts:

```go
New(cs.String(), WithMigrationsPath("./migrations"))
```

## Testing in Docker

See [here](https://github.com/alextanhongpin/postgres-test).
