package caching

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type IRequestStore interface {
	StoreRequest(*http.Request) (string, error)
	RetrieveRequest(string) (*CachedRequest, error)
}

type RequestStore struct {
	cache IServiceCache
}

func CreateRequestStore(cache IServiceCache) *RequestStore {
	return &RequestStore{cache: cache}
}

func (rs *RequestStore) StoreRequest(r *http.Request) (string, error) {
	requestID := uuid.New().String()
	cachedRequest, err := constructCachedRequest(r)
	if err != nil {
		return "", err
	}

	cachedBytes, err := json.Marshal(cachedRequest)
	if err != nil {
		return "", err
	}

	rs.cache.RequestAdd(requestID, cachedBytes)

	return requestID, nil
}

func (rs *RequestStore) RetrieveRequest(requestID string) (*CachedRequest, error) {
	cachedBytes, err := rs.cache.RequestRetrieve(requestID)
	if err != nil {
		return nil, err
	}

	var cachedRequest *CachedRequest
	err = json.Unmarshal(cachedBytes, cachedRequest)
	if err != nil {
		return nil, err
	}

	return cachedRequest, nil
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
