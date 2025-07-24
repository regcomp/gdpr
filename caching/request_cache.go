package caching

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

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
