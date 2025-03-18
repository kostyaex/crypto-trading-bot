package strategy

import (
	"encoding/json"
)

// Strategy представляет торговую стратегию
type Strategy struct {
	ID          int             `db:"id"`
	Name        string          `db:"name"`
	Description string          `db:"description"`
	Config      json.RawMessage `db:"config"`
	Active      bool            `db:"active"`
}

// NewStrategy создает новую торговую стратегию
func NewStrategy(name, description string, config map[string]interface{}) (*Strategy, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return &Strategy{
		Name:        name,
		Description: description,
		Config:      configJSON,
		Active:      true,
	}, nil
}

// UnmarshalConfig распаковывает конфигурацию стратегии
func (s *Strategy) UnmarshalConfig() (map[string]interface{}, error) {
	var config map[string]interface{}
	if err := json.Unmarshal(s.Config, &config); err != nil {
		return nil, err
	}
	return config, nil
}

// MarshalConfig упаковывает конфигурацию стратегии
func (s *Strategy) MarshalConfig(config map[string]interface{}) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	s.Config = configJSON
	return nil
}
