package settings

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type SettingsRegistry struct {
	factories map[string]func() Settings
	mu        sync.RWMutex
}

func NewSettingsRegistry() *SettingsRegistry {
	return &SettingsRegistry{
		factories: make(map[string]func() Settings),
	}
}

// Регистрация фабрики для типа компонента
func (r *SettingsRegistry) Register(settingsType string, factory func() Settings) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.factories[settingsType]; exists {
		panic(fmt.Sprintf("settings type %q already registered", settingsType))
	}
	r.factories[settingsType] = factory
}

// Создание и десериализация компонента по типу и сырым данным
func (r *SettingsRegistry) Build(settingsType string, rawJSON json.RawMessage) (Settings, error) {
	r.mu.RLock()
	factory, exists := r.factories[settingsType]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown component type: %s", settingsType)
	}

	instance := factory() // создаём пустой экземпляр структуры

	if err := json.Unmarshal(rawJSON, instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal component %s: %w", settingsType, err)
	}

	if err := validate.Struct(instance); err != nil {
		return nil, fmt.Errorf("validation failed for component %s: %w", settingsType, err)
	}

	return instance, nil
}
