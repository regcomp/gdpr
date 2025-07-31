package main

import "testing"

func TestGenerateFlattenedDataMapping(t *testing.T) {
	tcs := []struct {
		name     string
		input    ConfigData
		expected map[string]string
	}{
		{
			name: "routers",
			input: ConfigData{
				Routers: map[string]RouterConfig{
					"API": {
						PathPrefix: "/api",
						Endpoints: []string{
							"users",
							"items",
						},
					},
				},
			},
			expected: map[string]string{
				"ROUTERS.API.ENDPOINTS.ITEMS": "/api/items",
				"ROUTERS.API.ENDPOINTS.USERS": "/api/users",
				"ROUTERS.API.PATH_PREFIX":     "/api",
			},
		},
		{
			name: "service workers",
			input: ConfigData{
				ServiceWorkers: map[string]ServiceWorkerConfig{
					"CACHE": {
						Path:  "/worker/path.js",
						Scope: "/worker/scope",
					},
				},
			},
			expected: map[string]string{
				"SERVICE_WORKERS.CACHE.PATH":  "/worker/path.js",
				"SERVICE_WORKERS.CACHE.SCOPE": "/worker/scope",
			},
		},
		{
			name: "cookies",
			input: ConfigData{
				Cookies: map[string]string{
					"SESSION_ID": "session-id",
				},
			},
			expected: map[string]string{
				"COOKIES.SESSION_ID": "session-id",
			},
		},
		{
			name: "query parameters",
			input: ConfigData{
				QueryParameters: map[string]string{
					"REDIRECT_URL": "redirect-url",
				},
			},
			expected: map[string]string{
				"QUERY_PARAMETERS.REDIRECT_URL": "redirect-url",
			},
		},
		{
			name: "headers",
			input: ConfigData{
				Headers: map[string]string{
					"CUSTOM_HEADER": "Custom-Header",
				},
			},
			expected: map[string]string{
				"HEADERS.CUSTOM_HEADER": "Custom-Header",
			},
		},
		{
			name: "values",
			input: ConfigData{
				Values: map[string]string{
					"VALUE": "value",
				},
			},
			expected: map[string]string{
				"VALUES.VALUE": "value",
			},
		},
		{
			name: "config keys",
			input: ConfigData{
				ConfigKeys: map[string]string{
					"CONFIG_KEY": "CONFIG_KEY",
				},
			},
			expected: map[string]string{
				"CONFIG_KEYS.CONFIG_KEY": "CONFIG_KEY",
			},
		},
		{
			name: "form values",
			input: ConfigData{
				FormValues: map[string]string{
					"NONCE": "nonce",
				},
			},
			expected: map[string]string{
				"FORM_VALUES.NONCE": "nonce",
			},
		},
		{
			name: "local files",
			input: ConfigData{
				LocalFiles: map[string]string{
					"ENV": ".env",
				},
			},
			expected: map[string]string{
				"LOCAL_FILES.ENV": ".env",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := generateFlattenedDataMapping(&tc.input)
			if len(tc.expected) != len(result) {
				t.Errorf("expected length=%d !=  result length=%d", len(tc.expected), len(result))
				return
			}
			for key, val := range result {
				if _, ok := tc.expected[key]; !ok {
					t.Errorf("key=%s malformed", key)
					return
				}
				if tc.expected[key] != val {
					t.Errorf("value mismatch expected=%s, result=%s", tc.expected[key], val)
					return
				}
			}
		})
	}
}

func newFixtureConfigDataNoShared() *ConfigData {
	return &ConfigData{
		Routers: map[string]RouterConfig{
			"BASE": {
				PathPrefix: "/",
				Endpoints: []string{
					"healthz",
				},
			},
			"API": {
				PathPrefix: "/api",
				Endpoints: []string{
					"users",
					"items",
				},
			},
		},
		ServiceWorkers: map[string]ServiceWorkerConfig{
			"CACHE": {
				Path:  "/static/js/sw/cache.js",
				Scope: "/static",
			},
		},
		Cookies: map[string]string{
			"SESSION_ID": "session-id",
		},
		QueryParameters: map[string]string{
			"REDIRECT_URL": "redirect-url",
		},
		Headers: map[string]string{
			"CUSTOM_HEADER": "Custom-Header",
		},
		Values: map[string]string{
			"VALUE": "value",
		},
		RequestContextKeys: map[string]string{
			"CLAIMS": "claims",
		},
		ConfigKeys: map[string]string{
			"CONFIG_KEY": "CONFIG_KEY",
		},
		FormValues: map[string]string{
			"NONCE": "nonce",
		},
		LocalFiles: map[string]string{
			"ENV": ".env",
		},
	}
}

func newFixtureSharedData() map[string]map[string]string {
	return map[string]map[string]string{
		"TEST_SHARED": {
			"HEALTHZ_PATH":           "@ROUTERS.BASE.HEALTHZ",
			"API_USERS":              "@ROUTERS.API.USERS",
			"CACHE_WORKER_PATH":      "@SERVICE_WORKERS.CACHE.PATH",
			"CACHE_WORKER_SCOPE":     "@SERVICE_WORKERS.CACHE.SCOPE",
			"SESSION_ID_COOKIE_NAME": "@COOKIES.SESSION_ID",
		},
	}
}
