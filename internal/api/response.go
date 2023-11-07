package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"

	"shorturl/internal/urlservice"
)

// errInvalidRequest is returned when server should respond with http.StatusBadRequest or codes.InvalidArgument
var errInvalidRequest = errors.New("request contains invalid data")

func writeNotAllowed(w http.ResponseWriter, allowedMethods []string) {
	allowHeaderValue := strings.Join(allowedMethods, ", ")
	w.Header().Add("Allow", allowHeaderValue)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// writeResponse sets status code and writes JSON body depending on the request error
func writeResponse(w http.ResponseWriter, requestedURL string, requestHandlingError error) {
	w.Header().Add("Content-Type", "application/json")

	if requestHandlingError != nil {
		writeError(w, requestHandlingError)
		return
	}

	writeResult(w, requestedURL)
}

func writeError(w http.ResponseWriter, requestHandlingError error) {
	statusCode, _ := errorStatusCodes(requestHandlingError)
	w.WriteHeader(statusCode)

	respBody := struct {
		Error string `json:"error"`
	}{requestHandlingError.Error()}

	writeBody(w, respBody)
}

func writeResult(w http.ResponseWriter, result string) {
	respBody := struct {
		URL string `json:"url"`
	}{result}

	writeBody(w, respBody)
}

// errorStatusCodes handles not-nil request handling error and returns both HTTP
// and gRPC error codes that should be set in response.
func errorStatusCodes(requestHandlingError error) (httpCode int, gRPCCode codes.Code) {
	switch {
	case errors.Is(requestHandlingError, errInvalidRequest):
		return http.StatusBadRequest, codes.InvalidArgument
	case errors.Is(requestHandlingError, urlservice.ErrURLNotFound):
		return http.StatusNotFound, codes.NotFound
	default:
		return http.StatusInternalServerError, codes.Internal
	}
}

func writeBody(w http.ResponseWriter, respBody any) {
	resp, err := json.Marshal(respBody)
	if err != nil {
		logError("failed to marshal resp body", err)
		return
	}

	if _, err := w.Write(resp); err != nil {
		logError("failed to write resp body", err)
	}
}

func logError(msg string, err error) {
	slog.Error(msg, slog.String("error", err.Error()))
}
