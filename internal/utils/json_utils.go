package utils

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// StructToJSON преобразует структуру в JSON строку с использованием mapstructure
func StructToJSON(input interface{}) (string, error) {
	var m map[string]interface{}
	err := mapstructure.Decode(input, &m)
	if err != nil {
		return "", fmt.Errorf("failed to decode struct to map: %w", err)
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("failed to marshal map to JSON: %w", err)
	}

	return string(jsonData), nil
}
