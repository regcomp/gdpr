package main

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func toCamelCase(s string) string {
	titleCaser := cases.Title(language.English)
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == '.'
	})

	var result strings.Builder
	for _, part := range parts {
		result.WriteString(titleCaser.String(part))
	}
	return result.String()
}

func processRouters(data *ConfigData, td *ServiceTemplateData) {
	for name, config := range data.Routers {
		td.RouterPrefixes = append(td.RouterPrefixes, RouterPrefix{
			Name:   name,
			Value:  config.PathPrefix,
			GoName: fmt.Sprintf("Router%sPathPrefix", toCamelCase(name)),
		})
	}
	sort.Slice(td.RouterPrefixes, func(i, j int) bool {
		return td.RouterPrefixes[i].Name < td.RouterPrefixes[j].Name
	})
}

func processPaths(data *ConfigData, td *ServiceTemplateData) {
	for routerName, config := range data.Routers {
		for _, endpoint := range config.Endpoints {
			endpointName := strings.TrimPrefix(endpoint, "/")
			td.Endpoints = append(td.Endpoints, Endpoint{
				Name:     endpointName,
				Value:    endpoint,
				GoName:   fmt.Sprintf("Endpoint%s", toCamelCase(endpointName)),
				Category: strings.ToLower(routerName),
			})

			fullPath := config.PathPrefix
			if config.PathPrefix == "/" {
				fullPath = endpoint
			} else {
				if !strings.HasSuffix(config.PathPrefix, "/") {
					fullPath += "/"
				}
				fullPath += strings.TrimPrefix(endpoint, "/")
			}

			td.FullPaths = append(td.FullPaths, FullPath{
				Name:     fmt.Sprintf("%s_%s", routerName, endpointName),
				Value:    fullPath,
				GoName:   fmt.Sprintf("Path%s%s", toCamelCase(routerName), toCamelCase(endpointName)),
				Category: strings.ToLower(routerName),
			})
		}
	}
	sort.Slice(td.Endpoints, func(i, j int) bool {
		if td.Endpoints[i].Category == td.Endpoints[j].Category {
			return td.Endpoints[i].Name < td.Endpoints[j].Name
		}
		return td.Endpoints[i].Category < td.Endpoints[j].Category
	})
	sort.Slice(td.FullPaths, func(i, j int) bool {
		if td.FullPaths[i].Category == td.FullPaths[j].Category {
			return td.FullPaths[i].Name < td.FullPaths[j].Name
		}
		return td.FullPaths[i].Category < td.FullPaths[j].Category
	})
}

func processServiceWorkers(data *ConfigData, td *ServiceTemplateData) {
	for name, config := range data.ServiceWorkers {
		baseName := toCamelCase(name)
		td.ServiceWorkers = append(td.ServiceWorkers,
			ServiceWorker{
				Name:   name,
				Field:  "PATH",
				Value:  config.Path,
				GoName: fmt.Sprintf("Worker%sPath", baseName),
			},
			ServiceWorker{
				Name:   name,
				Field:  "SCOPE",
				Value:  config.Scope,
				GoName: fmt.Sprintf("Worker%sScope", baseName),
			},
		)
	}
}

func processConfigKeys(data *ConfigData, td *ServiceTemplateData) {
	var configSliceItems []string
	for _, key := range data.ConfigKeys {
		goName := fmt.Sprintf("Config%sKey", toCamelCase(key))
		td.ConfigKeys = append(td.ConfigKeys, ConfigKey{
			Name:   key,
			Value:  key,
			GoName: goName,
		})
		configSliceItems = append(configSliceItems, goName)
	}
	td.ConfigSlice = strings.Join(configSliceItems, ",\n\t")
}

func processCookies(data *ConfigData, td *ServiceTemplateData) {
	for name, cookie := range data.Cookies {
		goName := fmt.Sprintf("%sCookieName", toCamelCase(name))
		td.Cookies = append(td.Cookies, Cookie{
			Name:   name,
			Value:  cookie,
			GoName: goName,
		})
	}
}

func processRequestContextKeys(data *ConfigData, td *ServiceTemplateData) {
	for _, key := range data.RequestContextKeys {
		name := toCamelCase(strings.ReplaceAll(key, "-", "_"))
		td.ContextKeys = append(td.ContextKeys, ContextKey{
			Name:   name,
			Value:  key,
			GoName: fmt.Sprintf("ContextKey%s", name),
		})
	}
}

func processFormValues(data *ConfigData, td *ServiceTemplateData) {
	for _, value := range data.FormValues {
		name := toCamelCase(value)
		td.FormValues = append(td.FormValues, FormValue{
			Name:   name,
			Value:  value,
			GoName: fmt.Sprintf("FormValue%s", name),
		})
	}
}

func processLocalFiles(data *ConfigData, td *ServiceTemplateData) {
	for name, file := range data.LocalFiles {
		goName := fmt.Sprintf("Local%sPath", toCamelCase(name))
		td.LocalFiles = append(td.LocalFiles, LocalFile{
			Name:   name,
			Value:  file,
			GoName: goName,
		})
	}
}

func processHeaders(data *ConfigData, td *ServiceTemplateData) {
	for name, value := range data.Headers {
		goName := fmt.Sprintf("Header%s", toCamelCase(name))
		td.Headers = append(td.Headers, Header{
			Name:   name,
			Value:  value,
			GoName: goName,
		})
	}
}

func processValues(data *ConfigData, td *ServiceTemplateData) {
	for name, value := range data.Values {
		goName := fmt.Sprintf("Value%s", toCamelCase(name))
		td.Values = append(td.Values, Value{
			Name:   name,
			Value:  value,
			GoName: goName,
		})
	}
}

func mapRouterData(data *ConfigData, valueLookup map[string]string) {
	for name, config := range data.Routers {
		valueLookup[fmt.Sprintf("ROUTERS.%s.PATH_PREFIX", name)] = config.PathPrefix
		for _, endpoint := range config.Endpoints {
			fullPath := config.PathPrefix
			if config.PathPrefix == "/" {
				fullPath = endpoint
			} else {
				if !strings.HasSuffix(config.PathPrefix, "/") {
					fullPath += "/"
				}
				fullPath += strings.TrimPrefix(endpoint, "/")
			}
			// Create multiple key formats to handle different reference styles
			endpointKey := strings.TrimPrefix(endpoint, "/")

			// Original endpoint name (e.g., "renew-token")
			valueLookup[fmt.Sprintf("ROUTERS.%s.ENDPOINTS.%s", name, endpointKey)] = fullPath

			// Uppercase with underscores (e.g., "RENEW_TOKEN")
			endpointKeyUpper := strings.ToUpper(strings.ReplaceAll(endpointKey, "-", "_"))
			valueLookup[fmt.Sprintf("ROUTERS.%s.ENDPOINTS.%s", name, endpointKeyUpper)] = fullPath

			// Lowercase with underscores (e.g., "renew_token")
			endpointKeyLower := strings.ToLower(strings.ReplaceAll(endpointKey, "-", "_"))
			valueLookup[fmt.Sprintf("ROUTERS.%s.ENDPOINTS.%s", name, endpointKeyLower)] = fullPath
		}
	}
}

func mapHeaderData(data *ConfigData, valueLookup map[string]string) {
	for name, value := range data.Headers {
		valueLookup[fmt.Sprintf("HEADERS.%s", name)] = value
	}
}

func mapValuesData(data *ConfigData, valueLookup map[string]string) {
	for name, value := range data.Values {
		valueLookup[fmt.Sprintf("VALUES.%s", name)] = value
	}
}

func mapServiceWorkerData(data *ConfigData, valueLookup map[string]string) {
	for name, config := range data.ServiceWorkers {
		valueLookup[fmt.Sprintf("SERVICE_WORKERS.%s.PATH", name)] = config.Path
		valueLookup[fmt.Sprintf("SERVICE_WORKERS.%s.SCOPE", name)] = config.Scope
	}
}

func mapCookieData(data *ConfigData, valueLookup map[string]string) {
	for _, cookie := range data.Cookies {
		valueLookup[fmt.Sprintf("COOKIES.%s", strings.ToUpper(strings.ReplaceAll(cookie, "-", "_")))] = cookie
	}
}

func mapQueryParamData(data *ConfigData, valueLookup map[string]string) {
	for _, param := range data.QueryParameters {
		valueLookup[fmt.Sprintf("QUERY_PARAMETERS.%s", strings.ToUpper(strings.ReplaceAll(param, "-", "_")))] = param
	}
}
