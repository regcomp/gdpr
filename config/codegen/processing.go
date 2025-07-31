package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

func processData(data *ConfigData) *ServiceTemplateData {
	td := &ServiceTemplateData{}
	td.Header = WarningHeader

	processFuncs := []func(*ConfigData, *ServiceTemplateData){
		processRouters,
		processEndpoints,
		processPaths,
		processServiceWorkers,
		processConfigKeys,
		processCookies,
		processQueryParameters,
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

func processEndpoints(data *ConfigData, td *ServiceTemplateData) {
	for routerName, config := range data.Routers {
		for _, endpoint := range config.Endpoints {
			td.Endpoints = append(td.Endpoints, Endpoint{
				Name:     endpoint,
				Value:    fmt.Sprintf("/%s", endpoint),
				GoName:   fmt.Sprintf("Endpoint%s", toCamelCase(endpoint)),
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
}

func processPaths(data *ConfigData, td *ServiceTemplateData) {
	for routerName, config := range data.Routers {
		for _, endpoint := range config.Endpoints {
			td.FullPaths = append(td.FullPaths, FullPath{
				Name:     fmt.Sprintf("%s_%s", routerName, toKeyCase(endpoint)),
				Value:    path.Join(config.PathPrefix, endpoint),
				GoName:   fmt.Sprintf("Path%s%s", toCamelCase(routerName), toCamelCase(endpoint)),
				Category: strings.ToLower(routerName),
			})
		}
	}
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
	for _, key := range data.ConfigKeys {
		goName := fmt.Sprintf("Config%sKey", toCamelCase(key))
		td.ConfigKeys = append(td.ConfigKeys, ConfigKey{
			Name:   key,
			Value:  key,
			GoName: goName,
		})
		td.ConfigAttrs = append(td.ConfigAttrs, ConfigAttr{
			GoName: goName,
		})
	}
}

func processCookies(data *ConfigData, td *ServiceTemplateData) {
	for name, cookie := range data.Cookies {
		goName := fmt.Sprintf("CookieName%s", toCamelCase(name))
		td.Cookies = append(td.Cookies, Cookie{
			Name:   name,
			Value:  cookie,
			GoName: goName,
		})
	}
}

func processQueryParameters(data *ConfigData, td *ServiceTemplateData) {
	for name, parameter := range data.QueryParameters {
		goName := fmt.Sprintf("QueryParam%s", toCamelCase(name))
		td.QueryParams = append(td.QueryParams, QueryParam{
			Name:   name,
			Value:  parameter,
			GoName: goName,
		})
	}
}

func processRequestContextKeys(data *ConfigData, td *ServiceTemplateData) {
	for name, key := range data.RequestContextKeys {
		goName := fmt.Sprintf("ContextKey%s", toCamelCase(name))
		td.ContextKeys = append(td.ContextKeys, ContextKey{
			Name:   name,
			Value:  key,
			GoName: goName,
		})
	}
}

func processFormValues(data *ConfigData, td *ServiceTemplateData) {
	for name, value := range data.FormValues {
		goName := fmt.Sprintf("FormValue%s", toCamelCase(name))
		td.FormValues = append(td.FormValues, FormValue{
			Name:   name,
			Value:  value,
			GoName: goName,
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
