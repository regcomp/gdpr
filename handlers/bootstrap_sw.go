package handlers

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/regcomp/gdpr/templates/components"
)

func BootstrapServiceWorker(swPath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: find a place to put data models so this and the component can reference it
		req := struct {
			URL    string
			Method string
			Header map[string][]string
			Body   string
		}{
			URL:    r.URL.String(),
			Method: r.Method,
			Header: r.Header,
			Body:   extractBody(r),
		}

		w.Header().Add("Content-Type", "application/javascript")

		components.BootstrapServiceWorker(swPath, req.URL, req.Method, req.Body, req.Header).Render(r.Context(), w)
	})
}

func extractBody(r *http.Request) string {
	if r.Method == http.MethodGet || r.Method == http.MethodHead {
		return ""
	}

	contentType := r.Header.Get("Content-Type")
	body, _ := io.ReadAll(r.Body)

	// TODO: re-evaluate this logic
	if strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "text/") ||
		strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return string(body)
	} else {
		return base64.StdEncoding.EncodeToString(body)
	}
}
