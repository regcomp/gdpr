package main

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

func processData(data *ConfigData) (*ServiceTemplateData, error) {
	td := &ServiceTemplateData{}
	td.Header = WarningHeader

	processFuncs := []func(*ConfigData, *ServiceTemplateData) error{
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

	var err error
	for _, processFunc := range processFuncs {
		err = processFunc(data, td)
		if err != nil {
			return nil, err
		}
	}

	return td, nil
}

func resolveFullPaths(data *ConfigData) (map[string]string, error) {
	fullPaths := make(map[string]string)
	visited := make(map[string]bool)

	var resolvePath func(name string) (string, error)
	resolvePath = func(name string) (string, error) {
		if path, exists := fullPaths[name]; exists {
			return path, nil
		}

		if visited[name] {
			return "", fmt.Errorf("circular dependency detected involving router: %s", name)
		}

		config, exists := data.Routers[name]
		if !exists {
			return "", fmt.Errorf("router not found: %s", name)
		}

		visited[name] = true
		defer func() { visited[name] = false }()

		if config.Parent == nil {
			// Root router
			fullPaths[name] = config.PathPrefix
			return config.PathPrefix, nil
		}

		parentPath, err := resolvePath(*config.Parent)
		if err != nil {
			return "", err
		}

		fullPath := path.Join(parentPath, config.PathPrefix)
		fullPaths[name] = fullPath
		return fullPath, nil
	}

	// Resolve all paths
	for name := range data.Routers {
		if _, err := resolvePath(name); err != nil {
			return nil, err
		}
	}

	return fullPaths, nil
}

func processRouters(data *ConfigData, td *ServiceTemplateData) error {
	fullPaths, err := resolveFullPaths(data)
	if err != nil {
		return fmt.Errorf("failed to resolve router paths: %w", err)
	}

	for name, config := range data.Routers {
		parentName := ""
		if config.Parent != nil {
			parentName = *config.Parent
		}

		td.RouterPrefixes = append(td.RouterPrefixes, RouterPrefix{
			Name:      name,
			Value:     config.PathPrefix,
			FullValue: fullPaths[name],
			GoName:    fmt.Sprintf("Router%sPathPrefix", toCamelCase(name)),
			Parent:    parentName,
		})
	}

	sort.Slice(td.RouterPrefixes, func(i, j int) bool {
		return td.RouterPrefixes[i].Name < td.RouterPrefixes[j].Name
	})

	return nil
}

func processEndpoints(data *ConfigData, td *ServiceTemplateData) error {
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

	return nil
}

func processPaths(data *ConfigData, td *ServiceTemplateData) error {
	fullPaths, err := resolveFullPaths(data)
	if err != nil {
		return fmt.Errorf("failed to resolve router paths: %w", err)
	}

	for routerName, config := range data.Routers {
		for _, endpoint := range config.Endpoints {
			fullMountPath := fullPaths[routerName]

			td.FullPaths = append(td.FullPaths, FullPath{
				Name:     fmt.Sprintf("%s_%s", routerName, toKeyCase(endpoint)),
				Value:    path.Join(fullMountPath, endpoint),
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

	return nil
}

func processServiceWorkers(data *ConfigData, td *ServiceTemplateData) error {
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

	return nil
}

func processConfigKeys(data *ConfigData, td *ServiceTemplateData) error {
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

	return nil
}

func processCookies(data *ConfigData, td *ServiceTemplateData) error {
	for name, cookie := range data.Cookies {
		goName := fmt.Sprintf("CookieName%s", toCamelCase(name))
		td.Cookies = append(td.Cookies, Cookie{
			Name:   name,
			Value:  cookie,
			GoName: goName,
		})
	}

	return nil
}

func processQueryParameters(data *ConfigData, td *ServiceTemplateData) error {
	for name, parameter := range data.QueryParameters {
		goName := fmt.Sprintf("QueryParam%s", toCamelCase(name))
		td.QueryParams = append(td.QueryParams, QueryParam{
			Name:   name,
			Value:  parameter,
			GoName: goName,
		})
	}

	return nil
}

func processRequestContextKeys(data *ConfigData, td *ServiceTemplateData) error {
	for name, key := range data.RequestContextKeys {
		goName := fmt.Sprintf("ContextKey%s", toCamelCase(name))
		td.ContextKeys = append(td.ContextKeys, ContextKey{
			Name:   name,
			Value:  key,
			GoName: goName,
		})
	}

	return nil
}

func processFormValues(data *ConfigData, td *ServiceTemplateData) error {
	for name, value := range data.FormValues {
		goName := fmt.Sprintf("FormValue%s", toCamelCase(name))
		td.FormValues = append(td.FormValues, FormValue{
			Name:   name,
			Value:  value,
			GoName: goName,
		})
	}

	return nil
}

func processLocalFiles(data *ConfigData, td *ServiceTemplateData) error {
	for name, file := range data.LocalFiles {
		goName := fmt.Sprintf("Local%sPath", toCamelCase(name))
		td.LocalFiles = append(td.LocalFiles, LocalFile{
			Name:   name,
			Value:  file,
			GoName: goName,
		})
	}

	return nil
}

func processHeaders(data *ConfigData, td *ServiceTemplateData) error {
	for name, value := range data.Headers {
		goName := fmt.Sprintf("Header%s", toCamelCase(name))
		td.Headers = append(td.Headers, Header{
			Name:   name,
			Value:  value,
			GoName: goName,
		})
	}

	return nil
}

func processValues(data *ConfigData, td *ServiceTemplateData) error {
	for name, value := range data.Values {
		goName := fmt.Sprintf("Value%s", toCamelCase(name))
		td.Values = append(td.Values, Value{
			Name:   name,
			Value:  value,
			GoName: goName,
		})
	}

	return nil
}
