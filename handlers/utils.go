package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/google/uuid"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code >= 500 {
		log.Printf("Responding with 5xx error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func NewURL(scheme, host, path string) url.URL {
	return url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
}

type IRequestStore interface {
	StoreCachedRequest(r *http.Request) (string, error)
	RetrieveCachedRequest(string) (*CachedRequest, error)
}

type RequestStore struct {
	Store sync.Map
}

func CreateRequestStore() *RequestStore {
	return &RequestStore{Store: sync.Map{}}
}

// WARN: this struct with its field names are coupled with the script in /static/pages/register_service_worker.templ
type CachedRequest struct {
	URL    string              `json:"url"`
	Method string              `json:"method"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

func constructCachedRequest(r *http.Request) (CachedRequest, error) {
	body, err := extractBody(r)
	if err != nil {
		// TODO:
	}
	return CachedRequest{
		URL:    r.URL.String(),
		Method: r.Method,
		Header: r.Header.Clone(),
		Body:   body,
	}, nil
}

func extractBody(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO:
	}
	r.Body.Close()

	return string(body), nil
}

func (rs *RequestStore) StoreCachedRequest(r *http.Request) (string, error) {
	requestID := uuid.New().String()
	cachedValue, err := constructCachedRequest(r)
	if err != nil {
		return "", fmt.Errorf("RequestStore::StoreCachedRequest::could not cache request for redirect")
	}

	for {
		_, isCollision := rs.Store.LoadOrStore(requestID, cachedValue)
		if !isCollision {
			break
		}
		requestID = uuid.New().String()
	}

	return requestID, nil
}

func (rs *RequestStore) RetrieveCachedRequest(requestID string) (*CachedRequest, error) {
	cachedValue, exists := rs.Store.LoadAndDelete(requestID)
	if !exists {
		return nil, fmt.Errorf("RequestStore::GetCachedRequest::cached request does not exist")
	}
	cachedRequest, ok := cachedValue.(CachedRequest)
	if !ok {
		return nil, fmt.Errorf("RequestStore::GetCachedRequest::malformed cached request")
	}
	return &cachedRequest, nil
}
