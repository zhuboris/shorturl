package api

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"shorturl/internal/urlservice"
)

type requestResult struct {
	URL   any `json:"url"`
	Error any `json:"error"`
}

func TestRequestWithNotAllowedMethod(t *testing.T) {
	urlServiceMock := NewMockshortURLService(t)
	listenAddr := ":" + strconv.Itoa(rand.Intn(1e4))
	sut := NewRESTServer(listenAddr, urlServiceMock)

	request := httptest.NewRequest(http.MethodPatch, "/", nil)
	recorder := httptest.NewRecorder()
	sut.server.Handler.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
	assertBodyContent(t, recorder)
	require.NotEmpty(t, recorder.Header().Get("Allow"), "Allow header should be set")
	assert.Contains(t, recorder.Header().Get("Allow"), http.MethodPost)
	assert.Contains(t, recorder.Header().Get("Allow"), http.MethodGet)
}

func TestGetRequest(t *testing.T) {
	existingShortURL := "1234567890"

	tests := []struct {
		name               string
		path               string
		expectedStatusCode int
	}{
		{
			name:               "short url exists",
			path:               "/" + existingShortURL,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "short url does not exist",
			path:               "/1111111111",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "empty short url",
			path:               "/",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlServiceMock := NewMockshortURLService(t)
			if tt.expectedStatusCode != http.StatusBadRequest {
				urlServiceMock.EXPECT().
					OriginalURL(mock.Anything, mock.Anything).
					RunAndReturn(func(_ context.Context, shortURL string) (string, error) {
						if shortURL == existingShortURL {
							return "https://example.com/", nil
						}

						return "", urlservice.ErrURLNotFound
					}).
					Once()
			}

			listenAddr := ":" + strconv.Itoa(rand.Intn(1e4))
			sut := NewRESTServer(listenAddr, urlServiceMock)

			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			recorder := httptest.NewRecorder()
			sut.server.Handler.ServeHTTP(recorder, request)

			require.Equal(t, tt.expectedStatusCode, recorder.Code)
			assertBodyContent(t, recorder)
		})
	}
}

func TestPostRequest(t *testing.T) {
	tests := []struct {
		name               string
		originalURL        string
		expectedStatusCode int
	}{
		{
			name:               "valid url in request",
			originalURL:        "https://example.com/",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "empty request body",
			originalURL:        "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "invalid url in request body",
			originalURL:        "https://examp  le.com/",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlServiceMock := NewMockshortURLService(t)
			if tt.expectedStatusCode == http.StatusOK {
				urlServiceMock.EXPECT().
					ShortURL(mock.Anything, mock.Anything).
					Return("1111111111", nil).
					Once()
			}

			listenAddr := ":" + strconv.Itoa(rand.Intn(1e4))
			sut := NewRESTServer(listenAddr, urlServiceMock)

			bodyRaw := struct {
				URL string `json:"url"`
			}{tt.originalURL}
			requestBody, err := json.Marshal(bodyRaw)
			require.NoError(t, err, "Failed to marshal request body")

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBody))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()
			sut.server.Handler.ServeHTTP(recorder, request)

			require.Equal(t, tt.expectedStatusCode, recorder.Code)
			result := assertBodyContent(t, recorder)
			if recorder.Code != http.StatusOK {
				return
			}

			shortURL, ok := result.URL.(string)
			require.True(t, ok, "URL in response must be string")
			_, err = url.Parse(shortURL)
			assert.NoError(t, err, "URL in request is invalid")
		})
	}
}

func assertBodyContent(t *testing.T, recorder *httptest.ResponseRecorder) requestResult {
	t.Helper()

	var result requestResult
	if len(recorder.Body.Bytes()) != 0 {
		err := json.Unmarshal(recorder.Body.Bytes(), &result)
		require.NoError(t, err, "Error unmarshal resp body")
	}

	switch recorder.Code {
	case http.StatusOK:
		assert.NotEmpty(t, result.URL)
	case http.StatusMethodNotAllowed:
		assert.Empty(t, result)
	default:
		assert.NotEmpty(t, result.Error)
	}

	return result
}
