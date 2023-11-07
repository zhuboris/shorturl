package api

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"shorturl/internal/pb"
	"shorturl/internal/urlservice"
)

func TestGetOriginalURLMethod(t *testing.T) {
	existingShortURL := "1234567890"

	tests := []struct {
		name         string
		shortURL     string
		expectedCode codes.Code
	}{
		{
			name:         "short url exists",
			shortURL:     existingShortURL,
			expectedCode: codes.OK,
		},
		{
			name:         "short url does not exist",
			shortURL:     "1111111111",
			expectedCode: codes.NotFound,
		},
		{
			name:         "empty short url",
			shortURL:     "",
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlServiceMock := NewMockshortURLService(t)
			if tt.expectedCode != codes.InvalidArgument {
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

			client := grpcClient(t, urlServiceMock)
			originalURL, err := client.GetOriginalURL(context.Background(), &pb.ShortURL{Url: tt.shortURL})
			assertCorrectGRPCCode(t, err, tt.expectedCode)
			if status.Code(err) != codes.OK {
				return
			}

			assert.NotEmpty(t, originalURL.Url, "Missing response url")
		})
	}
}

func TestCreateShortURLMethod(t *testing.T) {
	tests := []struct {
		name         string
		originalURL  string
		expectedCode codes.Code
	}{
		{
			name:         "valid url in request",
			originalURL:  "https://example.com/",
			expectedCode: codes.OK,
		},
		{
			name:         "empty request body",
			originalURL:  "",
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "invalid url in request body",
			originalURL:  "https://examp  le.com/",
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlServiceMock := NewMockshortURLService(t)
			if tt.expectedCode != codes.InvalidArgument {
				urlServiceMock.EXPECT().
					ShortURL(mock.Anything, mock.Anything).
					Return("1111111111", nil).
					Once()
			}

			client := grpcClient(t, urlServiceMock)
			shortURL, err := client.CreateShortURL(context.Background(), &pb.OriginalURL{Url: tt.originalURL})
			assertCorrectGRPCCode(t, err, tt.expectedCode)
			if status.Code(err) != codes.OK {
				return
			}

			assert.NotEmpty(t, shortURL.Url, "Missing response url")
		})
	}
}

func assertCorrectGRPCCode(t *testing.T, err error, expectedCode codes.Code) {
	respStatus, _ := status.FromError(err)
	require.Equal(t, expectedCode, respStatus.Code())
	if respStatus.Code() == codes.OK {
		return
	}

	assert.NotEmpty(t, respStatus.Message(), "Error message must be set")
}

func grpcClient(t *testing.T, urlService ShortURLService) pb.ShortURLServiceClient {
	const bufSize = 1 << 20

	t.Helper()
	listener := bufconn.Listen(bufSize)
	runTestGRPCServer(t, urlService, listener)
	return connectGRPCClient(t, listener)
}

func connectGRPCClient(t *testing.T, listener *bufconn.Listener) pb.ShortURLServiceClient {
	t.Helper()
	dialOptionFunc := func(_ context.Context, _ string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(dialOptionFunc), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to make grpc dial connection")
	t.Cleanup(func() {
		conn.Close()
	})

	return pb.NewShortURLServiceClient(conn)
}

func runTestGRPCServer(t *testing.T, urlService ShortURLService, listener *bufconn.Listener) {
	t.Helper()

	serv := initGRPCServer(urlService)
	serv.listener = listener
	go func() {
		err := serv.Run()
		require.NoError(t, err)
	}()

	t.Cleanup(serv.Stop)
}
