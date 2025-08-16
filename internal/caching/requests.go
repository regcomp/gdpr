package caching

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type RequestManager struct {
	cache IRequestStash
}

func CreateRequestManager(cache IRequestStash) *RequestManager {
	return &RequestManager{cache: cache}
}

func (rs *RequestManager) StashRequest(r *http.Request) (string, error) {
	id, data, err := requestToCachedRequestAsBytes(r)
	if err != nil {
		return "", err
	}
	err = rs.cache.StashRequest(id, data)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (rs *RequestManager) RetrieveRequest(id string) (*CachedRequest, error) {
	data, err := rs.cache.RetrieveRequest(id)
	if err != nil {
		return nil, err
	}
	cachedRequest, err := bytesToCachedRequest(data)
	if err != nil {
		return nil, err
	}

	return cachedRequest, nil
}

// CachedRequest field names are coupled with the script in /static/pages/register_service_worker.templ
type CachedRequest struct {
	URL    string              `json:"url"`
	Method string              `json:"method"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

func requestToCachedRequestAsBytes(r *http.Request) (string, []byte, error) {
	requestID := uuid.New().String()
	cachedRequest, err := constructCachedRequest(r)
	if err != nil {
		return "", nil, err
	}

	cachedBytes, err := json.Marshal(cachedRequest)
	if err != nil {
		return "", nil, err
	}

	return requestID, cachedBytes, nil
}

func bytesToCachedRequest(data []byte) (*CachedRequest, error) {
	cachedRequest := &CachedRequest{}
	err := json.Unmarshal(data, cachedRequest)
	if err != nil {
		return nil, err
	}

	return cachedRequest, nil
}

func constructCachedRequest(r *http.Request) (CachedRequest, error) {
	body, err := extractBody(r)
	if err != nil {
		return CachedRequest{}, err
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
		return "", nil
	}
	r.Body.Close()

	return string(body), nil
}
