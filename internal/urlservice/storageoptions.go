package urlservice

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"shorturl/internal/encoder"
	"shorturl/internal/urlservice/dbstore"
	"shorturl/internal/urlservice/memstore"
)

// StorageOptionFunc is used to select a storage type for InMemoryURLStorage instance.
type StorageOptionFunc func(idEncoder encoder.IDEncoder, shortURLLength uint) urlStorage

// WithInMemoryStorage returns an option that initializes and returns in-memory storage for urlStorage interface.
func WithInMemoryStorage() StorageOptionFunc {
	return func(idEncoder encoder.IDEncoder, shortURLLength uint) urlStorage {
		inMemoryStorage := memstore.NewInMemoryURLStorage(idEncoder, shortURLLength)
		adapter := newInMemoryURLStorageAdapter(inMemoryStorage)

		return adapter
	}
}

// WithPostgreSQLStorage returns an option that initializes and returns PostgreSQL storage for urlStorage interface.
// It needs a connection pool to initialize database storage.
func WithPostgreSQLStorage(pool *pgxpool.Pool) StorageOptionFunc {
	return func(idEncoder encoder.IDEncoder, shortURLLength uint) urlStorage {
		return dbstore.NewPostgreSQLStorage(pool, idEncoder, shortURLLength)
	}
}
