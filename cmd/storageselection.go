package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"shorturl/internal/urlservice"
)

// selectedStorageOption is parsing flags and returning selected urlservice.StorageOptionFunc (or default).
//
// If a selected option does not exist, it returns error.
func selectedStorageOption() (urlservice.StorageOptionFunc, error) {
	const (
		inMemoryOption = "in-memory"
		postgresOption = "postgres"
	)

	tooltip := fmt.Sprintf("Specify the type of storage to use ('%s' or '%s'). Default is '%s'", inMemoryOption, postgresOption, inMemoryOption)
	storageType := flag.String("s", inMemoryOption, tooltip)
	flag.Parse()

	switch *storageType {
	case inMemoryOption:
		return urlservice.WithInMemoryStorage(), nil
	case postgresOption:
		return withPostgresStorage()
	default:
		return nil, fmt.Errorf("invalid input: got %q, valid options: %q, %q", inMemoryOption, postgresOption, *storageType)
	}
}

func withPostgresStorage() (urlservice.StorageOptionFunc, error) {
	pool, err := postgresPool()
	if err != nil {
		return nil, err
	}

	return urlservice.WithPostgreSQLStorage(pool), nil
}

// postgresPool initializes postgres connection pool with values from environment variables.
// It has a timeout for connection and returns error on connection fails.
func postgresPool() (*pgxpool.Pool, error) {
	const timeoutValue = 15 * time.Second

	var (
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		host     = os.Getenv("POSTGRES_HOST")
		port     = os.Getenv("POSTGRES_PORT")
		dbName   = os.Getenv("POSTGRES_DB")
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeoutValue)
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("faild to connect to postgres: %w", err)
	}

	return pool, nil
}
