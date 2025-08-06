package config

import (
	"fmt"
	"log"
	"strings"
)

type IConfigStore interface {
	GetServiceURL() string
	GetServiceURLWithPort() string
	GetDefaultPort() string
	GetSessionDuration() string

	GetSecretStoreConfig() *SecretStoreConfig
	GetServiceCacheConfig() *ServiceCacheConfig
	GetAuthProviderConfig() *AuthProviderConfig
	GetDatabaseManagerConfig() (*DatabaseManagerConfig, error)
}

type LocalConfigStore struct {
	// mu    sync.RWMutex // May need in the future
	attrs map[string]string
}

func NewConfigStore(getenv func(string) string, getters ...func(string) string) IConfigStore {
	storeType := getenv(ConfigConfigStoreTypeKey)
	switch storeType {
	case ValueLocalType:
		return newLocalConfigStore(getters...)
	default:
		log.Fatalf("unknown config store type: %s", storeType)
		return nil
	}
}

func newLocalConfigStore(getters ...func(string) string) *LocalConfigStore {
	store := &LocalConfigStore{
		attrs: make(map[string]string),
	}
	store.initializeLocalStore(getters...)
	return store
}

func (cs *LocalConfigStore) initializeLocalStore(getters ...func(string) string) {
	for _, getter := range getters {
		for _, attr := range ConfigAttrs {
			cs.attrs[attr] = getter(attr)
		}
	}
}

func (cs *LocalConfigStore) GetServiceURL() string {
	return cs.attrs[ConfigServiceUrlKey]
}

func (cs *LocalConfigStore) GetDefaultPort() string {
	return cs.attrs[ConfigDefaultPortKey]
}

func (cs *LocalConfigStore) GetServiceURLWithPort() string {
	return cs.attrs[ConfigServiceUrlKey] + ":" + cs.GetDefaultPort()
}

func (cs *LocalConfigStore) GetSessionDuration() string {
	return cs.attrs[ConfigSessionDurationKey]
}

func (cs *LocalConfigStore) GetSecretStoreConfig() *SecretStoreConfig {
	return &SecretStoreConfig{
		StoreType: cs.attrs[ConfigSecretStoreTypeKey],
	}
}

func (cs *LocalConfigStore) GetServiceCacheConfig() *ServiceCacheConfig {
	return &ServiceCacheConfig{
		CacheType: cs.attrs[ConfigServiceCacheTypeKey],
	}
}

func (cs *LocalConfigStore) GetAuthProviderConfig() *AuthProviderConfig {
	return &AuthProviderConfig{
		ProviderType: cs.attrs[ConfigAuthProviderTypeKey],
	}
}

func (cs *LocalConfigStore) GetDatabaseManagerConfig() (*DatabaseManagerConfig, error) {
	recordsConfig, err := cs.createRecordsConfig()
	if err != nil {
		return nil, err
	}

	providerConfigs, errs := cs.createProviderConfigs()
	if errs != nil {
		return nil, err
	}

	return &DatabaseManagerConfig{
		RecordsConfig:   recordsConfig,
		ProviderConfigs: providerConfigs,
	}, nil
}

func (cs *LocalConfigStore) createRecordsConfig() (RecordsDatabaseConfig, error) {
	recordsType, ok := cs.attrs[ConfigRecordsDatabaseTypeKey]
	if !ok {
		return RecordsDatabaseConfig{}, fmt.Errorf("no config entry found for records database")
	}
	return RecordsDatabaseConfig{ProviderType: recordsType}, nil
}

func (cs *LocalConfigStore) createProviderConfigs() ([]ProviderDatabaseConfig, error) {
	providerNamesRaw, ok := cs.attrs[ConfigDatabaseProviderNamesKey]
	if !ok {
		return nil, fmt.Errorf("no config entry found for provider names")
	}
	names := parseProviderNames(providerNamesRaw)

	providerTypesRaw, ok := cs.attrs[ConfigDatabaseProviderTypesKey]
	if !ok {
		return nil, fmt.Errorf("no config entry found for provider types")
	}
	nameToType := parseProviderTypes(providerTypesRaw)

	providerTableNamesRaw, ok := cs.attrs[ConfigDatabaseProviderTableNamesKey]
	if !ok {
		return nil, fmt.Errorf("no config entry found for table names")
	}
	nameToTableNamesArray := parseProviderTableNames(providerTableNamesRaw)

	providerConfigs := make([]ProviderDatabaseConfig, 0)

	if len(names) != len(nameToType) || len(names) != len(nameToTableNamesArray) {
		return nil, fmt.Errorf("malformed config. langth of names=%d, types=%d, tables=%d should be equal",
			len(names), len(nameToType), len(nameToTableNamesArray))
	}

	for _, name := range names {
		cfg := ProviderDatabaseConfig{
			ProviderName: name,
			ProviderType: nameToType[name],
			TableNames:   nameToTableNamesArray[name],
		}
		providerConfigs = append(providerConfigs, cfg)
	}

	return providerConfigs, nil
}

// parseProviderNames example:"db1,db2,db3"
func parseProviderNames(names string) []string {
	return parseItems(names)
}

// parseProviderTypes example:"db1:AWSS3;db2:PLANETSCALE"
func parseProviderTypes(cfgString string) map[string]string {
	entries := parseEntries(cfgString)
	typesMap := make(map[string]string, 0)
	for _, entry := range entries {
		pName, pType := parseByName(entry)
		typesMap[pName] = pType
	}
	return typesMap
}

// parseProviderTableNames example:"db1:table1,table2,table3;db2:table1,table2"
func parseProviderTableNames(tableNames string) map[string][]string {
	entries := parseEntries(tableNames)
	result := make(map[string][]string, 0)
	for _, entry := range entries {
		pName, tNamesString := parseByName(entry)
		tNames := parseItems(tNamesString)
		result[pName] = tNames
	}

	return result
}

func parseEntries(s string) []string {
	return strings.Split(s, ValueEntryDelim)
}

func parseByName(s string) (name, value string) {
	result := strings.SplitN(s, ValueNameDelim, 2)
	return result[0], result[1]
}

func parseItems(s string) []string {
	return strings.Split(s, ValueItemDelim)
}
