// Package urlservice provides a service that creates and manages URL storage.
package urlservice

import (
	"context"
	"errors"
	"fmt"

	"shorturl/internal/encoder"
)

type urlStorage interface {
	OriginalURL(ctx context.Context, shortURL string) (string, error)
	ShortURL(ctx context.Context, originalURL string) (string, error)
}

// ShortURLService is a service to manipulate with selected URL storage.
//
// It must be initialized with NewShortURLService to set desired storage.
type ShortURLService struct {
	storage urlStorage
}

// NewShortURLService initializes a new ShortURLService instance with a storage, chosen with StorageOptionFunc.
// It also takes id encoder and length of short URL to set up storage.
func NewShortURLService(idEncoder encoder.IDEncoder, shortURLLength uint, storageOption StorageOptionFunc) ShortURLService {
	storage := storageOption(idEncoder, shortURLLength)
	return ShortURLService{storage}
}

// ErrURLNotFound is returned when provided short URL not maps with any original URL.
var ErrURLNotFound = errors.New("requested short url has no matches")

// OriginalURL calls method OriginalURL in his storage and returns ErrURLNotFound if
// the method returned an error.
func (s ShortURLService) OriginalURL(ctx context.Context, shortURL string) (string, error) {
	original, err := s.storage.OriginalURL(ctx, shortURL)
	if err != nil {
		return "", errors.Join(ErrURLNotFound, err)
	}

	return original, nil
}

// ShortURL calls method ShortURL in his storage.
func (s ShortURLService) ShortURL(ctx context.Context, originalURL string) (string, error) {
	short, err := s.storage.ShortURL(ctx, originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to insert or get short url for url %q: %w", originalURL, err)
	}

	return short, nil
}
