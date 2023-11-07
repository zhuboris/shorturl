package urlservice

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"shorturl/internal/urlservice/memstore"
)

type encoderStub struct{}

func (e encoderStub) EncodeID(_, _ uint) string {
	return ""
}

func Test_newInMemoryURLStorageAdapter(t *testing.T) {
	storage := memstore.NewInMemoryURLStorage(encoderStub{}, 10)
	adapter := newInMemoryURLStorageAdapter(storage)

	assert.Implements(t, (*urlStorage)(nil), adapter)
}
