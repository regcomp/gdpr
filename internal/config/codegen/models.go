package main

type RouterConfig struct {
	PathPrefix string   `json:"PATH_PREFIX"`
	Endpoints  []string `json:"ENDPOINTS"`
	Parent     *string  `json:"PARENT"`
}

type ServiceWorkerConfig struct {
	Path  string `json:"PATH"`
	Scope string `json:"SCOPE"`
}

type ConfigData struct {
	Routers            map[string]RouterConfig        `json:"ROUTERS"`
	ServiceWorkers     map[string]ServiceWorkerConfig `json:"SERVICE_WORKERS"`
	Cookies            map[string]string              `json:"COOKIES"`
	QueryParameters    map[string]string              `json:"QUERY_PARAMETERS"`
	Headers            map[string]string              `json:"HEADERS"`
	Values             map[string]string              `json:"VALUES"`
	RequestContextKeys map[string]string              `json:"REQUEST_CONTEXT_KEYS"`
	ConfigKeys         map[string]string              `json:"CONFIG_KEYS"`
	FormValues         map[string]string              `json:"FORM_VALUES"`
	LocalFiles         map[string]string              `json:"LOCAL_FILES"`
	Shared             map[string]map[string]string   `json:"SHARED"`
}

type ServiceTemplateData struct {
	Header         string
	RouterPrefixes []RouterPrefix
	Endpoints      []Endpoint
	FullPaths      []FullPath
	ServiceWorkers []ServiceWorker
	ConfigKeys     []ConfigKey
	ConfigAttrs    []ConfigAttr // used to iterate over
	Cookies        []Cookie
	QueryParams    []QueryParam
	ContextKeys    []ContextKey
	FormValues     []FormValue
	LocalFiles     []LocalFile
	Headers        []Header
	Values         []Value
	SharedValues   []SharedValue
}

type RouterPrefix struct {
	Name      string
	Value     string
	FullValue string
	GoName    string
	Parent    string
}

type Endpoint struct {
	Name     string
	Value    string
	GoName   string
	Category string // used for comments in the template and sorting
}

type FullPath struct {
	Name     string
	Value    string
	GoName   string
	Category string // used for comments in the template and sorting
}

type ServiceWorker struct {
	Name   string
	Field  string // which field of the struct this represents
	Value  string
	GoName string
}

type ConfigKey struct {
	Name   string
	Value  string
	GoName string
}

type ConfigAttr struct {
	GoName string
}

type Cookie struct {
	Name   string
	Value  string
	GoName string
}

type QueryParam struct {
	Name   string
	Value  string
	GoName string
}

type ContextKey struct {
	Name   string
	Value  string
	GoName string
}

type FormValue struct {
	Name   string
	Value  string
	GoName string
}

type LocalFile struct {
	Name   string
	Value  string
	GoName string
}

type Header struct {
	Name   string
	Value  string
	GoName string
}

type Value struct {
	Name   string
	Value  string
	GoName string
}

type SharedValue struct {
	Name   string
	Value  string
	GoName string
}
