package urlservice

import (
	"context"

	"shorturl/internal/urlservice/memstore"
)

// inMemoryURLStorageAdapter is used as type InMemoryURLStorage to implement the urlStorage interface.
type inMemoryURLStorageAdapter struct {
	storage *memstore.InMemoryURLStorage
}

func newInMemoryURLStorageAdapter(storage *memstore.InMemoryURLStorage) inMemoryURLStorageAdapter {
	return inMemoryURLStorageAdapter{storage}
}

func (a inMemoryURLStorageAdapter) OriginalURL(_ context.Context, shortURL string) (string, error) {
	return a.storage.OriginalURL(shortURL)
}

func (a inMemoryURLStorageAdapter) ShortURL(_ context.Context, originalURL string) (string, error) {
	return a.storage.ShortURL(originalURL)
}
