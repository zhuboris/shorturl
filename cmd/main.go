// Package main is the entry point for shorturl_api program. It initializes URL
// storage and runs web servers. To select a storage type on program launch, it
// supports a flag "-s [option]":
//   - option for in-memory storage is "in-memory"
//   - option for PostgreSQL storage is "postgres"
//
// The default option is in-memory.
//
// These environment variables must be set for work:
//   - "HTTP_LISTEN_ADDRESS": listen address for REST API server
//   - "GRPC_LISTEN_ADDRESS": listen address for gRPC server
//
// If selected PostgreSQL storage, these variables must also be set:
//   - "POSTGRES_HOST": host of postgres server
//   - "POSTGRES_PORT": port of postgres server
//   - "POSTGRES_USER": name of postgres user
//   - "POSTGRES_PASSWORD": password of that user
//   - "POSTGRES_DB": name of database
//
// Also, it supports optional "SHORT_URL_LENGTH" variable to set desired length
// of short URL, the default value is 10.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"shorturl/internal/api"
	"shorturl/internal/encoder"
	"shorturl/internal/urlservice"
)

func main() {
	err := run()
	log.Fatal("Program is shutdown", err)
}

// run is a function to start program and return its possible errors. It gets
// selected options, initializes servers and starts them. Also, it processes
// shutdown on reading a first message from server's error channel. This message
// means that some server is down.
func run() error {
	idEncoder := encoder.NewIDEncoder()
	shortURLService, err := initShortURLService(idEncoder)
	if err != nil {
		return err
	}

	gRPCServer, restServer, err := initServers(shortURLService)
	if err != nil {
		return err
	}

	errCh := runServers(restServer, gRPCServer)
	defer shutdownServers(restServer, gRPCServer)

	err = <-errCh
	shutdownError := shutdownServers(restServer, gRPCServer)

	return errors.Join(err, shutdownError)
}

func initShortURLService(idEncoder encoder.IDEncoder) (urlservice.ShortURLService, error) {
	storageOption, err := selectedStorageOption()
	if err != nil {
		return urlservice.ShortURLService{}, err
	}

	shortURLLength, err := lookForShortURLLength()
	if err != nil {
		return urlservice.ShortURLService{}, err
	}

	shortURLService := urlservice.NewShortURLService(idEncoder, uint(shortURLLength), storageOption)
	return shortURLService, nil
}

func initServers(shortURLService urlservice.ShortURLService) (*api.GRPCServer, *api.RESTServer, error) {
	gRPCAddress := os.Getenv("GRPC_LISTEN_ADDRESS")
	gRPCServer, err := api.NewGRPCServer(gRPCAddress, shortURLService)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init gRPC server: %w", err)
	}

	restAPIAddress := os.Getenv("HTTP_LISTEN_ADDRESS")
	restServer := api.NewRESTServer(restAPIAddress, shortURLService)
	return gRPCServer, restServer, nil
}

// runServers starts both servers in goroutines and writes their result errors to chanel.
// It returns the read-only channel to get messages about shutdown of servers.
func runServers(restServer *api.RESTServer, gRPCServer *api.GRPCServer) <-chan error {
	errCh := make(chan error)
	go func() {
		errCh <- restServer.Run()
	}()

	go func() {
		errCh <- gRPCServer.Run()
	}()

	return errCh
}

func shutdownServers(restServer *api.RESTServer, gRPCServer *api.GRPCServer) error {
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	gRPCServer.Stop()
	return restServer.Shutdown(ctx)
}

func lookForShortURLLength() (int, error) {
	const defaultLength = 10

	raw, isSet := os.LookupEnv("SHORT_URL_LENGTH")
	if !isSet {
		return defaultLength, nil
	}

	result, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("short url env contains not int: %w", err)
	}

	return result, nil
}
