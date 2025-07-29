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
	WarningHeader = "// Auto-generated from /config/service.json - DO NOT EDIT\n"
)

const (
	sourceServiceFile             = "../../config/service.json"
	targetGoFilePath              = "../../constants/service.go"
	javascriptSharedDirectoryPath = "../../static/js/shared/"
	jsTemplatePath                = "./templates/js.gotmpl"
	goTemplatePath                = "./templates/go.gotmpl"
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

	templateData := processData(data)

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

	if err := os.MkdirAll(javascriptSharedDirectoryPath, 0755); err != nil {
		log.Fatal("Error creating output directory:", err)
	}

	for subfieldName, subfieldData := range data.Shared {
		filename := strings.ToLower(subfieldName) + "_shared.js"
		filepath := fmt.Sprintf("%s/%s", javascriptSharedDirectoryPath, filename)

		sharedData := generateSharedDataMapping(data, subfieldData)

		generateSingleJSFile(filepath, subfieldName, sharedData)
	}
}

func processData(data *ConfigData) *ServiceTemplateData {
	td := &ServiceTemplateData{}
	td.Header = WarningHeader

	processFuncs := []func(*ConfigData, *ServiceTemplateData){
		processRouters,
		processPaths,
		processServiceWorkers,
		processConfigKeys,
		processCookies,
		processRequestContextKeys,
		processFormValues,
		processLocalFiles,
		processHeaders,
		processValues,
	}

	for _, processFunc := range processFuncs {
		processFunc(data, td)
	}

	return td
}

func generateSharedDataMapping(data *ConfigData, subfieldMap map[string]string) map[string]string {
	sharedData := make(map[string]string)
	valueLookup := make(map[string]string)

	mappingFuncs := []func(*ConfigData, map[string]string){
		mapRouterData,
		mapHeaderData,
		mapValuesData,
		mapServiceWorkerData,
		mapCookieData,
		mapQueryParamData,
	}

	for _, mappingFunc := range mappingFuncs {
		mappingFunc(data, valueLookup)
	}

	// Process the subfield data
	for key, value := range subfieldMap {
		if after, ok := strings.CutPrefix(value, "@"); ok {
			refKey := after
			if refValue, exists := valueLookup[refKey]; exists {
				sharedData[key] = refValue
			} else {
				sharedData[key] = value
			}
		} else {
			sharedData[key] = value
		}
	}

	return sharedData
}

func generateSingleJSFile(filepath string, subfieldName string, resolvedData map[string]string) {
	constantName := strings.ToUpper(subfieldName) + "_CONSTANTS"

	templateData := struct {
		SubfieldName string
		ConstantName string
		Data         map[string]string
	}{
		SubfieldName: subfieldName,
		ConstantName: constantName,
		Data:         resolvedData,
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
