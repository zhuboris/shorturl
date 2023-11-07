// Package dbstore provides a database URL storage.
//
// It allows for the storing and retrieval of original URLs using
// encoded keys.
package dbstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"shorturl/internal/encoder"
)

// PostgreSQLStorage is a database URL storage using PostgreSQL.
//
// The zero value is not useful, you must use NewPostgreSQLStorage to create an instance.
type PostgreSQLStorage struct {
	pool           *pgxpool.Pool
	idEncoder      encoder.IDEncoder
	shortURLLength uint
}

// NewPostgreSQLStorage initializes a new PostgreSQLStorage instance with the given database connection pool,
// ID encoder, and the specified length for short URLs. It returns a pointer to created object.
func NewPostgreSQLStorage(pool *pgxpool.Pool, idEncoder encoder.IDEncoder, shortURLLength uint) *PostgreSQLStorage {
	return &PostgreSQLStorage{
		pool:           pool,
		idEncoder:      idEncoder,
		shortURLLength: shortURLLength,
	}
}

// OriginalURL is looking for the original URL by passed short URL.
// If the short URL does not exist in the storage, it returns an error.
func (s PostgreSQLStorage) OriginalURL(ctx context.Context, shortURL string) (string, error) {
	const sql = `
		SELECT original_url FROM short_urls
		WHERE url = $1;
	`

	var originalURL string
	err := s.pool.QueryRow(ctx, sql, shortURL).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("failed to get %q url from db: %w", shortURL, err)
	}

	return originalURL, nil
}

// ShortURL should always return short URL for provided original URL value.
// At first, it tries to find saved value, but if it does not exist, it encodes the
// original URL by incremented ID from the database and returns a new value.
//
// If inserting transaction returned error that some unique value is already saved in the database,
// the function checks once more if the original URL exists.
// It can happen in a case where another same transaction was executed before.
func (s PostgreSQLStorage) ShortURL(ctx context.Context, originalURL string) (string, error) {
	shortURL, err := s.tryFindShortURL(ctx, originalURL)
	if !errors.Is(err, pgx.ErrNoRows) {
		return shortURL, fmt.Errorf("unexpected eror in db for %q url: %w", originalURL, err)
	}

	var newSearchError error
	shortURL, insertError := s.addNewURL(ctx, originalURL, shortURL)
	if isURLAddedByOtherTransaction(insertError) {
		shortURL, newSearchError = s.tryFindShortURL(ctx, originalURL)
	}

	return shortURL, fmt.Errorf("failed to get %q url from db: %w: %w", originalURL, insertError, newSearchError)
}

func (s PostgreSQLStorage) tryFindShortURL(ctx context.Context, originalURL string) (string, error) {
	const sql = `
		SELECT url FROM short_urls
		WHERE original_url = $1;
	`

	var shortURL string
	err := s.pool.QueryRow(ctx, sql, originalURL).Scan(&shortURL)
	return shortURL, err
}

func (s PostgreSQLStorage) addNewURL(ctx context.Context, originalURL string, shortURL string) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}

	defer tx.Rollback(ctx)
	newID, err := s.insertOriginalURL(ctx, originalURL, tx)
	if err != nil {
		return "", err
	}

	shortURL, err = s.setShortURL(ctx, originalURL, newID, tx)
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	return shortURL, err
}

func isURLAddedByOtherTransaction(insertError error) bool {
	const postgresUniqueDuplicateErrorCode = "23505"

	var postgresError *pgconn.PgError
	if !errors.As(insertError, &postgresError) {
		return false
	}

	return postgresError.Code == postgresUniqueDuplicateErrorCode
}

func (s PostgreSQLStorage) insertOriginalURL(ctx context.Context, originalURL string, tx pgx.Tx) (uint, error) {
	const sql = `
		INSERT INTO original_urls (url)
		VALUES ($1)
		
		RETURNING id;
	`
	var newID uint
	err := tx.QueryRow(ctx, sql, originalURL).Scan(&newID)

	return newID, err
}

func (s PostgreSQLStorage) setShortURL(ctx context.Context, originalURL string, urlID uint, tx pgx.Tx) (string, error) {
	const sql = `
		INSERT INTO short_urls (original_url, url) 
		VALUES ($1, $2);
	`
	shortURL := s.idEncoder.EncodeID(urlID, s.shortURLLength)
	_, err := tx.Exec(ctx, sql, originalURL, shortURL)

	return shortURL, err
}
