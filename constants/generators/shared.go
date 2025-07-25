package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	jsonData, err := os.ReadFile("../../constants/shared.json")
	if err != nil {
		log.Fatal("Error reading headers.json:", err)
	}

	var sharedData struct {
		Headers map[string]string `json:"headers"`
	}

	// unmarshalling is for json validation
	if err := json.Unmarshal(jsonData, &sharedData); err != nil {
		log.Fatal("Invalid JSON:", err)
	}

	jsContent := fmt.Sprintf(`// Auto-generated from /constants/shared.json - DO NOT EDIT
const SHARED = %s;

// Export for modules if needed
if (typeof module !== 'undefined' && module.exports) {
    module.exports = HEADERS;
}
`, string(jsonData))

	// Write to static directory
	if err := os.MkdirAll("../../static", 0755); err != nil {
		log.Fatal("Error creating static directory:", err)
	}

	if err := os.WriteFile("../../static/js/shared.js", []byte(jsContent), 0644); err != nil {
		log.Fatal("Error writing headers.js:", err)
	}
}
