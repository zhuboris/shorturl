// Package memstore provides an in-memory URL storage.
//
// It allows for the storing and retrieval of original URLs using
// encoded keys. It's safe for concurrent use.
package memstore

import (
	"errors"
	"fmt"
	"sync"

	"shorturl/internal/encoder"
)

// InMemoryURLStorage is an in-memory storage for URLs.
//
// It maps both original URL by encoded URLs and encoded urls by original URLs.
// This is needed for fast search both values.
// Encoding depends on URL id, so it also stores the current value of incrementing id.
//
// The zero value is not useful, you must use NewInMemoryURLStorage to create an instance.
type InMemoryURLStorage struct {
	originalByEncodedURLs map[string]string
	encodedByOriginalURLs map[string]string
	idEncoder             encoder.IDEncoder
	currentID             uint
	shortURLLength        uint
	mutex                 sync.RWMutex
}

// NewInMemoryURLStorage initializes a new InMemoryURLStorage instance with the given ID encoder
// and the specified length for short URLs. It returns a pointer to created object.
func NewInMemoryURLStorage(idEncoder encoder.IDEncoder, shortURLLength uint) *InMemoryURLStorage {
	return &InMemoryURLStorage{
		idEncoder:             idEncoder,
		shortURLLength:        shortURLLength,
		encodedByOriginalURLs: make(map[string]string),
		originalByEncodedURLs: make(map[string]string),
	}
}

// OriginalURL is looking for the original URL by passed short URL.
// If the short URL does not exist in the storage, it returns an error.
func (s *InMemoryURLStorage) OriginalURL(shortURL string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	originalURL, isFound := s.originalByEncodedURLs[shortURL]
	if !isFound {
		return "", fmt.Errorf("%q url not found in im-memory storage", shortURL)
	}

	return originalURL, nil
}

// ShortURL should always return short URL for provided original URL value.
//
// At first, it tries to find saved value, but if it does not exist, it encodes the
// original URL by incremented ID and returns a new value.
//
// saveNewURL method call is tying to get saved value again before saving a new one
// due to a case where another goroutine has already performed this operation before.
//
// It might return an error if encoder returns a short URL that already exists in storage
// or if encoded value has an incorrect length.
func (s *InMemoryURLStorage) ShortURL(originalURL string) (string, error) {
	shortURL, isFound := s.lookForShortURL(originalURL)
	if isFound {
		return shortURL, nil
	}

	return s.saveNewURL(originalURL)
}

func (s *InMemoryURLStorage) lookForShortURL(originalURL string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	shortURL, isFound := s.encodedByOriginalURLs[originalURL]
	return shortURL, isFound
}

func (s *InMemoryURLStorage) saveNewURL(toAdd string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if shortURL, isAddedAlready := s.encodedByOriginalURLs[toAdd]; isAddedAlready {
		return shortURL, nil
	}

	s.currentID++
	newShortURL := s.idEncoder.EncodeID(s.currentID, s.shortURLLength)
	if err := s.checkResultLength(newShortURL); err != nil {
		return "", err
	}

	if err := s.checkIfResultUnique(newShortURL); err != nil {
		return "", err
	}

	s.originalByEncodedURLs[newShortURL] = toAdd
	s.encodedByOriginalURLs[toAdd] = newShortURL

	return newShortURL, nil
}

func (s *InMemoryURLStorage) checkIfResultUnique(shortURL string) error {
	if _, containsShortURL := s.originalByEncodedURLs[shortURL]; containsShortURL {
		return errors.New("encoded url is not unique")
	}

	return nil
}

func (s *InMemoryURLStorage) checkResultLength(shortURL string) error {
	if len(shortURL) != int(s.shortURLLength) {
		return fmt.Errorf("unexpected length of encoded url, expected=%d, actual=%d", s.shortURLLength, len(shortURL))
	}

	return nil
}
