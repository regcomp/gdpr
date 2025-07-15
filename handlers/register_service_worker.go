package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/regcomp/gdpr/templates/pages"
)

func (stx *ServiceContext) RegisterServiceWorker(swPath, swScope string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stx.RequestTracer.UpdateRequestTrace(r, "RegisterServiceWorker")
		req, err := constructRequestObject(r)
		if err != nil {
			log.Panicf("could not construct jsonRequest=%s", err.Error())
		}

		pages.RegisterServiceWorker(swPath, swScope, req).Render(r.Context(), w)
	})
}

// WARN: the json fields are coupled with bootstrapCodePath
type jsonRequest struct {
	URL    string              `json:"url"`
	Method string              `json:"method"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

func constructRequestObject(r *http.Request) (jsonRequest, error) {
	body, err := extractBody(r)
	if err != nil {
		// TODO:
	}
	return jsonRequest{
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
