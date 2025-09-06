package types

import "encoding/json"

// Component — общий интерфейс для всех конфигурируемых компонентов системы.
// Используется реестром для создания и возврата типизированных настроек.
type Component interface {
	// ComponentType возвращает уникальный строковый идентификатор типа компонента.
	// Должен совпадать с ID, по которому компонент регистрируется в реестре.
	ComponentType() string
}

type ComponentConfig struct {
	ID       string          `db:"id" json:"id"`
	Name     string          `db:"name" json:"name"`
	Settings json.RawMessage `db:"settings" json:"settings"`
}
