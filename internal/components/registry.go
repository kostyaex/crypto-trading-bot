package components

import (
	"crypto-trading-bot/internal/types"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ComponentRegistry struct {
	factories map[string]func() types.Component
	mu        sync.RWMutex
}

func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		factories: make(map[string]func() types.Component),
	}
}

// Регистрация фабрики для типа компонента
func (r *ComponentRegistry) Register(componentType string, factory func() types.Component) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.factories[componentType]; exists {
		panic(fmt.Sprintf("component type %q already registered", componentType))
	}
	r.factories[componentType] = factory
}

// Создание и десериализация компонента по типу и сырым данным
func (r *ComponentRegistry) Build(componentType string, rawJSON json.RawMessage) (types.Component, error) {
	r.mu.RLock()
	factory, exists := r.factories[componentType]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown component type: %s", componentType)
	}

	instance := factory() // создаём пустой экземпляр структуры

	if err := json.Unmarshal(rawJSON, instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal component %s: %w", componentType, err)
	}

	if err := validate.Struct(instance); err != nil {
		return nil, fmt.Errorf("validation failed for component %s: %w", componentType, err)
	}

	return instance, nil
}
