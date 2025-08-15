package config

import (
	"fmt"
	"strings"
)

type LocalConfigStore struct {
	// mu    sync.RWMutex // May need in the future
	attrs map[string]string
}

func newLocalConfigStore(getters ...func(string) string) (*LocalConfigStore, error) {
	store := &LocalConfigStore{
		attrs: make(map[string]string),
	}
	store.initializeLocalStore(getters...)
	return store, nil
}

func (cs *LocalConfigStore) initializeLocalStore(getters ...func(string) string) {
	for _, getter := range getters {
		for _, attr := range ConfigAttrs {
			cs.attrs[attr] = getter(attr)
		}
	}
}

func (cs *LocalConfigStore) GetServiceURL() (string, error) {
	serviceURL, ok := cs.attrs[ConfigServiceUrlKey]
	if !ok {
		return "", fmt.Errorf("could not find service url")
	}
	return serviceURL, nil
}

func (cs *LocalConfigStore) GetDefaultPort() (string, error) {
	defaultPort, ok := cs.attrs[ConfigDefaultPortKey]
	if !ok {
		return "", fmt.Errorf("could not find default port")
	}
	return defaultPort, nil
}

func (cs *LocalConfigStore) GetServiceURLWithPort() (string, error) {
	defaultPort, err := cs.GetDefaultPort()
	if err != nil {
		return "", err
	}
	serviceURL, err := cs.GetServiceURL()
	if err != nil {
		return "", err
	}
	return serviceURL + ":" + defaultPort, nil
}

func (cs *LocalConfigStore) GetSessionDuration() (string, error) {
	sessionDuration, ok := cs.attrs[ConfigSessionDurationKey]
	if !ok {
		return "", fmt.Errorf("could not find session duration")
	}

	return sessionDuration, nil
}

func (cs *LocalConfigStore) GetSecretStoreConfig() (*SecretStoreConfig, error) {
	storeType, okay := cs.attrs[ConfigSecretStoreTypeKey]
	if !okay {
		return nil, fmt.Errorf("could not find secret store type")
	}
	return &SecretStoreConfig{
		StoreType: storeType,
	}, nil
}

func (cs *LocalConfigStore) GetServiceCacheConfig() (*ServiceCacheConfig, error) {
	cacheType, ok := cs.attrs[ConfigServiceCacheTypeKey]
	if !ok {
		return nil, fmt.Errorf("could not find service cache type")
	}
	return &ServiceCacheConfig{
		CacheType: cacheType,
	}, nil
}

func (cs *LocalConfigStore) GetAuthProviderConfig() (*AuthProviderConfig, error) {
	providerType, ok := cs.attrs[ConfigAuthProviderTypeKey]
	if !ok {
		return nil, fmt.Errorf("could not find auth provider type")
	}
	return &AuthProviderConfig{
		ProviderType: providerType,
	}, nil
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
