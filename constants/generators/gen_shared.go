package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const generationHeader = "// Auto-generated from /constants/shared.json - DO NOT EDIT"

func main() {
	jsonData, err := os.ReadFile("../../constants/shared.json")
	if err != nil {
		log.Fatal("Error reading headers.json:", err)
	}

	var sharedData struct {
		Headers map[string]string `json:"HEADERS"`
		Paths   map[string]string `json:"PATHS"`
		Values  map[string]string `json:"VALUES"`
	}

	// unmarshalling is for json validation
	if err := json.Unmarshal(jsonData, &sharedData); err != nil {
		log.Fatal("Invalid JSON:", err)
	}

	jsContent := fmt.Sprintf(`%s
const SHARED = %s;`, generationHeader, string(jsonData))

	// Write to static directory
	if err := os.MkdirAll("../../static", 0755); err != nil {
		log.Fatal("Error creating static directory:", err)
	}

	if err := os.WriteFile("../../static/js/shared.js", []byte(jsContent), 0644); err != nil {
		log.Fatalf("Error writing /static/js/shared.js, err=%s", err)
	}

	goContent := fmt.Sprintf(`%s

package constants

const (
	// headers
	HeaderNonceToken             = "%s"
	HeaderRenewAccessToken       = "%s"
	HeaderServiceWorkerAllowed   = "%s"
	HeaderAuthRetryWorkerRunning = "%s"

	// paths
	PathAuthRenewToken = "%s"

	// string values
	ValueTrueString = "%s"
)

`, generationHeader,
		// headers
		sharedData.Headers["NONCE_TOKEN"],
		sharedData.Headers["RENEW_ACCESS_TOKEN"],
		sharedData.Headers["SERVICE_WORKER_ALLOWED"],
		sharedData.Headers["AUTH_RETRY_WORKER_RUNNING"],
		// paths
		sharedData.Paths["AUTH_RENEW"],
		// values
		sharedData.Values["TRUE"],
	)

	if err := os.WriteFile("../../constants/shared.go", []byte(goContent), 0644); err != nil {
		log.Fatalf("Error writing /constants/shared.go, err=%s", err)
	}
}
