// Package api provides implementation of both a gRPC and a REST servers.
// It handles incoming requests.
package api

import (
	"context"
	"errors"
	"fmt"
)

// ShortURLService is a definition of service that exchanges and stores URLs.
type ShortURLService interface {
	OriginalURL(ctx context.Context, shortURL string) (string, error)
	ShortURL(ctx context.Context, originalURL string) (string, error)
}

func handleCreationShortURL(ctx context.Context, originalURL string, urlService ShortURLService) (string, error) {
	parsedURL, err := validateURL(originalURL)
	if err != nil {
		return "", errors.Join(errInvalidRequest, err)
	}

	shortURL, err := urlService.ShortURL(ctx, parsedURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func handleGetOriginalURL(ctx context.Context, shortURL string, urlService ShortURLService) (string, error) {
	if shortURL == "" {
		return "", fmt.Errorf("%w: short url is not provided", errInvalidRequest)
	}

	return urlService.OriginalURL(ctx, shortURL)
}
