package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
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

func exemptServiceWorker(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/static/sw/") {
	}
}
