package engine

import (
	"crypto-trading-bot/pkg/types"
	"log"
	"sync"
)

var (
	registry   = make(map[string]ComponentFactory)
	registryMu sync.RWMutex
)

// RegisterComponent добавляет фабрику компонента в реестр
func RegisterComponent(componentType string, factory ComponentFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if _, exists := registry[componentType]; exists {
		log.Printf("Warning: component type '%s' is already registered, overwriting", componentType)
	}
	registry[componentType] = factory
}

// New создаёт компонент по типу
func New(
	componentType string,
	config types.ComponentConfig,
	logger *log.Logger,
) (*Component, error) {
	registryMu.RLock()
	factory, exists := registry[componentType]
	registryMu.RUnlock()

	if !exists {
		return nil, &UnknownComponentError{Type: componentType}
	}

	return factory(config, logger)
}

// ListRegistered возвращает список всех зарегистрированных типов (для дебага)
func ListRegistered() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	var types []string
	for k := range registry {
		types = append(types, k)
	}
	return types
}

// Ошибка: неизвестный тип компонента
type UnknownComponentError struct {
	Type string
}

func (e *UnknownComponentError) Error() string {
	return "unknown component type: " + e.Type
}
