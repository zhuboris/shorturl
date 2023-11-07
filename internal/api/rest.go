package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// RESTServer is REST API server implementation that processing requests to short URL service.
// It must be initialized with NewRESTServer
type RESTServer struct {
	server     *http.Server
	urlService ShortURLService
}

// NewRESTServer initializes RESTServer with its address to listen, and short URL service.
// It returns a pointer to object.
func NewRESTServer(listenAddress string, urlService ShortURLService) *RESTServer {
	server := &RESTServer{
		urlService: urlService,
	}

	server.initHTTPServer(listenAddress)
	return server
}

// Run is calling method ListenAndServe of object's http.Server and will return its error.
func (s *RESTServer) Run() error {
	return s.server.ListenAndServe()
}

// Shutdown is calling method Shutdown of object's http.Server and will return its error.
func (s *RESTServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// initHTTPServer is setting http.Server field of the object with its handler registration.
func (s *RESTServer) initHTTPServer(listenAddress string) {
	mux := http.NewServeMux()
	mux.Handle("/", loggingMiddleware(s.handleHTTP()))

	s.server = &http.Server{
		Handler: mux,
		Addr:    listenAddress,
	}
}

// handleHTTP is a handler for "/" path. It determines the request method and
// calls the corresponding handler. If the method is not allowed, it sets Allow
// handler and returns code 405. Else it processes the request and handles its result.
// It will write needed status code and body with a result url or error message.
func (s *RESTServer) handleHTTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resultURL string
		var err error
		switch r.Method {
		case http.MethodPost:
			resultURL, err = s.handlePost(r)
		case http.MethodGet:
			resultURL, err = s.handleGet(r)
		default:
			allowedMethods := []string{http.MethodPost, http.MethodGet}
			writeNotAllowed(w, allowedMethods)
			return
		}

		writeResponse(w, resultURL, err)
	}
}

func (s *RESTServer) handlePost(r *http.Request) (string, error) {
	rawURL, err := rawURLFromRequestBody(r)
	if err != nil {
		return "", errors.Join(errInvalidRequest, err)
	}

	return handleCreationShortURL(r.Context(), rawURL, s.urlService)
}

func (s *RESTServer) handleGet(r *http.Request) (string, error) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	return handleGetOriginalURL(r.Context(), shortURL, s.urlService)
}

func rawURLFromRequestBody(r *http.Request) (string, error) {
	var body struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("invalid json with url: %w", err)
	}

	return body.URL, nil
}

func validateURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", errors.New("missing url in request")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid url: %w", err)
	}

	return parsedURL.String(), nil
}
