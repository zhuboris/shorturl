package urlservice

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"shorturl/internal/encoder"
)

func TestNewShortURLService(t *testing.T) {
	optionMock := NewMockStorageOptionFunc(t)
	optionMock.EXPECT().
		Execute(mock.Anything, mock.Anything).
		RunAndReturn(func(_ encoder.IDEncoder, _ uint) urlStorage {
			mockStorage := NewMockurlStorage(t)
			return mockStorage
		})

	tests := []struct {
		name          string
		storageOption StorageOptionFunc
	}{
		{
			name:          "option mock",
			storageOption: optionMock.Execute,
		},
		{
			name:          "inMemoryStorage",
			storageOption: WithInMemoryStorage(),
		},
		{
			name:          "inMemoryStorage",
			storageOption: WithPostgreSQLStorage(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewShortURLService(encoderStub{}, 10, optionMock.Execute)
			assert.NotNil(t, service.storage)
		})
	}
}

func TestShortURLService_OriginalURL(t *testing.T) {
	tests := []struct {
		name          string
		wantError     bool
		expectedError error
	}{
		{
			name:      "no error",
			wantError: false,
		},
		{
			name:          "error",
			wantError:     true,
			expectedError: ErrURLNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockurlStorage(t)
			storageMock.EXPECT().
				OriginalURL(mock.Anything, mock.Anything).
				RunAndReturn(func(_ context.Context, _ string) (string, error) {
					if tt.wantError {
						return "", errors.New("some error")
					}

					return "", nil
				}).
				Once()

			sut := ShortURLService{
				storage: storageMock,
			}

			_, err := sut.OriginalURL(context.Background(), "123")
			if !tt.wantError {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}

func TestShortURLService_ShortURL(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "no error",
			wantError: false,
		},
		{
			name:      "error",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := NewMockurlStorage(t)
			storageMock.EXPECT().
				ShortURL(mock.Anything, mock.Anything).
				RunAndReturn(func(_ context.Context, _ string) (string, error) {
					if tt.wantError {
						return "", errors.New("some error")
					}

					return "", nil
				})

			sut := ShortURLService{
				storage: storageMock,
			}

			_, err := sut.ShortURL(context.Background(), "123")
			if !tt.wantError {
				assert.NoError(t, err)
				return
			}

			require.Error(t, err)
		})
	}
}
