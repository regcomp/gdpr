package main

import (
	"fmt"
	"path"
)

func generateFlattenedDataMapping(data *ConfigData) map[string]string {
	flattenedDataMapping := make(map[string]string)

	flatteningFuncs := []func(*ConfigData, map[string]string){
		flattenRouterData,
		flattenPathData,
		flattenServiceWorkerData,
		flattenCookieData,
		flattenQueryParamData,
		flattenConfigKeysData,
		// request context keys
		flattenFormValuesData,
		flattenLocalFilesData,
		flattenHeaderData,
		flattenValuesData,
	}

	for _, flatteningFunc := range flatteningFuncs {
		flatteningFunc(data, flattenedDataMapping)
	}

	return flattenedDataMapping
}

// NOTE: if the routing structure ever becomes deeper/more complex then this
// will need a different approach.
func flattenRouterData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.Routers {
		valueLookup[fmt.Sprintf("ROUTERS.%s.PATH_PREFIX", k)] = v.PathPrefix
	}
}

func flattenPathData(data *ConfigData, valueLookup map[string]string) {
	for routerName, routerConfig := range data.Routers {
		for _, endpoint := range routerConfig.Endpoints {
			fullPath := path.Join(routerConfig.PathPrefix, endpoint)
			endpointName := toKeyCase(endpoint, "-")
			constructedKey := fmt.Sprintf("ROUTERS.%s.ENDPOINTS.%s", routerName, endpointName)
			valueLookup[constructedKey] = fullPath
		}
	}
}

func flattenHeaderData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.Headers {
		valueLookup[fmt.Sprintf("HEADERS.%s", k)] = v
	}
}

func flattenValuesData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.Values {
		valueLookup[fmt.Sprintf("VALUES.%s", k)] = v
	}
}

func flattenServiceWorkerData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.ServiceWorkers {
		valueLookup[fmt.Sprintf("SERVICE_WORKERS.%s.PATH", k)] = v.Path
		valueLookup[fmt.Sprintf("SERVICE_WORKERS.%s.SCOPE", k)] = v.Scope
	}
}

func flattenCookieData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.Cookies {
		valueLookup[fmt.Sprintf("COOKIES.%s", k)] = v
	}
}

func flattenQueryParamData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.QueryParameters {
		valueLookup[fmt.Sprintf("QUERY_PARAMETERS.%s", k)] = v
	}
}

func flattenConfigKeysData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.ConfigKeys {
		valueLookup[fmt.Sprintf("CONFIG_KEYS.%s", k)] = v
	}
}

func flattenFormValuesData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.FormValues {
		valueLookup[fmt.Sprintf("FORM_VALUES.%s", k)] = v
	}
}

func flattenLocalFilesData(data *ConfigData, valueLookup map[string]string) {
	for k, v := range data.LocalFiles {
		valueLookup[fmt.Sprintf("LOCAL_FILES.%s", k)] = v
	}
}
