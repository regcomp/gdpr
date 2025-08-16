//go:generate go run ./

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	WarningHeader = "// Auto-generated from /config/config.json - DO NOT EDIT\n"
)

const (
	sourceServiceFile = "../../config/config.json"
	targetGoFilePath  = "../../config/constants.go"
	generatedTSPath   = "../../../web/generated"
	jsTemplatePath    = "./templates/js.gotmpl"
	goTemplatePath    = "./templates/go.gotmpl"
)

func main() {
	jsonData, err := os.ReadFile(sourceServiceFile)
	if err != nil {
		log.Fatal("Error reading config.json:", err)
	}

	data := &ConfigData{}
	if err := json.Unmarshal(jsonData, data); err != nil {
		log.Fatal("Invalid JSON:", err)
	}

	templateData, err := processData(data)
	if err != nil {
		log.Panic(err.Error())
	}

	tmpl, err := template.ParseFiles(goTemplatePath)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	goFile, err := os.Create(targetGoFilePath)
	if err != nil {
		log.Fatal("Error creating Go file:", err)
	}
	defer goFile.Close()

	if err := tmpl.Execute(goFile, templateData); err != nil {
		log.Fatal("Error executing template:", err)
	}

	if len(data.Shared) == 0 {
		return
	}

	if err := os.MkdirAll(generatedTSPath, 0755); err != nil {
		log.Fatal("Error creating output directory:", err)
	}

	// generating all the specified .js files to share constants
	flattenedMapping := generateFlattenedDataMapping(data)
	for subfieldName, subfieldData := range data.Shared {
		filename := strings.ToLower(subfieldName) + ".constants.ts"
		filepath := fmt.Sprintf("%s/%s", generatedTSPath, filename)

		sharedData := generateSharedDataMapping(flattenedMapping, subfieldData)

		generateSingleJSFile(filepath, subfieldName, sharedData)
	}
}

func generateSharedDataMapping(flattenedMapping map[string]string, subfieldData map[string]string) map[string]string {
	sharedData := make(map[string]string)

	for finalName, referenceKey := range subfieldData {
		trimmedKey := strings.TrimPrefix(referenceKey, "@")
		if referencedValue, exists := flattenedMapping[trimmedKey]; exists {
			sharedData[finalName] = referencedValue
		} else {
			sharedData[finalName] = "NULL"
		}
	}

	return sharedData
}

func generateSingleJSFile(filepath string, subfieldName string, resolvedData map[string]string) {
	dependentFilename := strings.ToLower(subfieldName)
	constantName := strings.ReplaceAll(strings.ToUpper(subfieldName), ".", "_") + "_CONSTANTS"

	// NOTE: these field names are coupled with the template
	templateData := struct {
		DependentFilename string
		ConstantName      string
		Data              map[string]string
	}{
		DependentFilename: dependentFilename,
		ConstantName:      constantName,
		Data:              resolvedData,
	}

	template, err := template.ParseFiles(jsTemplatePath)
	if err != nil {
		log.Fatal("Error parsing JavaScript template:", err)
	}

	jsFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal("Error creating JavaScript file:", err)
	}
	defer jsFile.Close()

	if err := template.Execute(jsFile, templateData); err != nil {
		log.Fatal("Error executing JavaScript template:", err)
	}
}
