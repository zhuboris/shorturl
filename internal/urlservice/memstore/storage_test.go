package memstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shorturl/internal/encoder"
)

type encoderStub struct{}

func (e encoderStub) EncodeID(_, _ uint) string {
	return "encoded"
}

var stubReturnValue = encoderStub{}.EncodeID(0, 0)

func TestNewInMemoryURLStorage(t *testing.T) {
	var shortURLLength uint = 10
	result := NewInMemoryURLStorage(encoderStub{}, shortURLLength)

	assert.Equal(t, shortURLLength, result.shortURLLength)
	assert.NotNil(t, result.originalByEncodedURLs, "Map originals was not init")
	assert.NotNil(t, result.encodedByOriginalURLs, "Map shorts was not init")
}

func TestInMemoryURLStorage_OriginalURL(t *testing.T) {
	tests := []struct {
		name           string
		originals      map[string]string
		shortURL       string
		expectedResult string
		requireError   require.ErrorAssertionFunc
	}{
		{
			name:           "short url exist",
			originals:      map[string]string{"short": "original"},
			shortURL:       "short",
			expectedResult: "original",
			requireError:   require.NoError,
		},
		{
			name:         "short url not exist",
			originals:    map[string]string{"short": "original"},
			shortURL:     "123",
			requireError: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewInMemoryURLStorage(encoderStub{}, 10)
			sut.originalByEncodedURLs = tt.originals

			result, err := sut.OriginalURL(tt.shortURL)
			tt.requireError(t, err)
			if err != nil {
				return
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestInMemoryURLStorage_ShortURL(t *testing.T) {
	tests := []struct {
		name           string
		shortURLLength int
		originals      map[string]string
		shorts         map[string]string
		originalURL    string
		expectedResult string
		requireError   require.ErrorAssertionFunc
	}{
		{
			name:           "original url exist",
			shortURLLength: len(stubReturnValue),
			originals:      map[string]string{"short": "original"},
			shorts:         map[string]string{"original": "short"},
			originalURL:    "original",
			expectedResult: "short",
			requireError:   require.NoError,
		},
		{
			name:           "original url not exist when short was not collided",
			shortURLLength: len(stubReturnValue),
			originals:      map[string]string{"short": "original"},
			shorts:         map[string]string{"original": "short"},
			originalURL:    "new",
			requireError:   require.NoError,
		},
		{
			name:           "original url not exist when short collided",
			shortURLLength: len(stubReturnValue),
			originals:      map[string]string{stubReturnValue: "original"},
			shorts:         map[string]string{"original": stubReturnValue},
			originalURL:    "new",
			requireError:   require.Error,
		},
		{
			name:           "length is greater than requested",
			shortURLLength: len(stubReturnValue) - 1,
			originals:      map[string]string{"short": "original"},
			shorts:         map[string]string{"original": "short"},
			originalURL:    "new",
			requireError:   require.Error,
		},
		{
			name:           "length is less than requested",
			shortURLLength: len(stubReturnValue) + 1,
			originals:      map[string]string{"short": "original"},
			shorts:         map[string]string{"original": "short"},
			originalURL:    "new",
			requireError:   require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewInMemoryURLStorage(encoderStub{}, uint(tt.shortURLLength))
			sut.originalByEncodedURLs = tt.originals
			sut.encodedByOriginalURLs = tt.shorts

			result, err := sut.ShortURL(tt.originalURL)
			tt.requireError(t, err)
			if err != nil {
				assert.NotContains(t, sut.originalByEncodedURLs, tt.originals, "URL should not be added on error")
				return
			}

			if tt.expectedResult != "" {
				assert.Equal(t, tt.expectedResult, result)
			}

			assert.Equal(t, sut.encodedByOriginalURLs[tt.originalURL], result, "Original url was not saved")
			assert.Equal(t, sut.originalByEncodedURLs[result], tt.originalURL, "Short url was not saved")
		})
	}
}

func TestInMemoryURLStorage_ShortURL_CheckThatIDIsIncrementing(t *testing.T) {
	idEncoder := encoder.NewIDEncoder()
	sut := NewInMemoryURLStorage(idEncoder, 10)
	result1, err := sut.ShortURL("url1")
	require.NoError(t, err)

	result2, err := sut.ShortURL("url2")
	require.NoError(t, err, "Short probably same")

	assert.NotEqual(t, result1, result2)
}

func TestInMemoryURLStorage_saveNewURL_ShouldReturnSavedResult_WhenResultWasAlreadyAdded(t *testing.T) {
	idEncoder := encoder.NewIDEncoder()
	sut := NewInMemoryURLStorage(idEncoder, 10)

	result1, err := sut.saveNewURL("url")
	require.NoError(t, err)

	result2, err := sut.saveNewURL("url")
	require.NoError(t, err)

	assert.Equal(t, result1, result2, "First result was not checked")
}
