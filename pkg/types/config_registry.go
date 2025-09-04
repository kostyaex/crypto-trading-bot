package types

import (
	"encoding/json"
	"fmt"
)

var configRegistry = map[string]func() ComponentConfig{
	"clustering": func() ComponentConfig { return &ClusterConfig{} },
	"dispatcher": func() ComponentConfig { return &DispatcherConfig{} },
}

// CreateConfig создаёт пустой экземпляр по типу
func CreateConfig(typ string) (ComponentConfig, error) {
	if factory, ok := configRegistry[typ]; ok {
		return factory(), nil
	}
	return nil, fmt.Errorf("unknown config type: %s", typ)
}

func DeserializeConfig(configType string, raw json.RawMessage) (ComponentConfig, error) {

	config, err := CreateConfig(configType)
	if err != nil {
		return nil, err
	}

	//json.Unmarshal знает тип через config
	if err := json.Unmarshal(raw, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config for %s: %w", configType, err)
	}

	return config, nil
}
