package models

import (
	"crypto-trading-bot/pkg/types"
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

// Strategy представляет торговую стратегию
type Strategy struct {
	ID          int             `db:"id"`
	Name        string          `db:"name"`
	Description string          `db:"description"`
	Config      json.RawMessage `db:"config"`
	Active      bool            `db:"active"`
	Settings    types.StrategyConfig
}

// NewStrategy создает новую торговую стратегию
// func NewStrategy(name, description string, settings StrategySettings) (*Strategy, error) {

// 	configJSON, err := utils.StructToJSON(settings)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Strategy{
// 		Name:                name,
// 		Description:         description,
// 		Config:              []byte(configJSON),
// 		Active:              true,
// 		SeriesBuilderConfig: settings.SeriesBuilderConfig,
// 	}, nil
// }

// Распаковывает JSON поле в структуру настроек
func (s *Strategy) UpdateSettingsFromConf() error {
	var config map[string]interface{}
	if err := json.Unmarshal(s.Config, &config); err != nil {
		return err
	}

	var settings types.StrategyConfig
	err := mapstructure.Decode(config, &settings)
	if err != nil {
		return err
	}

	s.Settings = settings

	return nil
}

// // метаданные поля настроек. Читается из тегов полей стркутур *Settings
// type SettingMetadata struct {
// 	Name        string
// 	Description string
// 	Nested      []SettingMetadata
// }

// func GetMetadataSettings(v interface{}) []SettingMetadata {
// 	meta := make([]SettingMetadata, 0)

// 	val := reflect.ValueOf(v)
// 	typ := val.Type()

// 	for i := 0; i < typ.NumField(); i++ {
// 		field := typ.Field(i)
// 		tag := field.Tag.Get("name")
// 		fmt.Printf("Field: %s, Tag: %s\n", field.Name, tag)
// 		metaGroup := SettingMetadata{
// 			Name:   tag,
// 			Nested: make([]SettingMetadata, 0),
// 		}

// 		// Если поле - вложенная структура, обходим её поля
// 		if field.Type.Kind() == reflect.Struct {
// 			//nestedVal := val.Field(i)
// 			nestedType := field.Type
// 			for j := 0; j < nestedType.NumField(); j++ {
// 				nestedField := nestedType.Field(j)
// 				nestedTag := nestedField.Tag.Get("name")
// 				fmt.Printf("  Nested Field: %s, Tag: %s\n", nestedField.Name, nestedTag)
// 				metaGroup.Nested = append(metaGroup.Nested, SettingMetadata{
// 					Name: nestedTag,
// 				})
// 			}
// 		}

// 		meta = append(meta, metaGroup)
// 	}

// 	return meta
// }

// func GetMetadataStrategyWavesSettings() []SettingMetadata {
// 	return GetMetadataSettings(StrategyWavesSettings{})
// }

// // Для веб-интерфейса. Выгружает JSON пустой структуры настроек.
// func StrategyWavesSettingsEmptyJson() string {
// 	//jsonData, err := json.Marshal(StrategyWavesSettings{})
// 	jsonData, err := utils.StructToJSON(StrategyWavesSettings{})
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal settings: %w", err).Error()
// 	}
// 	return string(jsonData)
// }

//func printTags(v interface{}) {
//val := reflect.ValueOf(v)
//typ := val.Type()

//for i := 0; i < typ.NumField(); i++ {
//field := typ.Field(i)
//tag := field.Tag.Get("name")
//fmt.Printf("Field: %s, Tag: %s\n", field.Name, tag)

//// Если поле - вложенная структура, обходим её поля
//if field.Type.Kind() == reflect.Struct {
////nestedVal := val.Field(i)
//nestedType := field.Type
//for j := 0; j < nestedType.NumField(); j++ {
//nestedField := nestedType.Field(j)
//nestedTag := nestedField.Tag.Get("name")
//fmt.Printf("  Nested Field: %s, Tag: %s\n", nestedField.Name, nestedTag)
//}
//}
//}
//}

// Функция для получения метаданных полей - извлекает метаданные из тегов структур и возвращает их в виде карты
//func getFieldMetadata(t interface{}) map[string]map[string]string {
//fieldMeta := make(map[string]map[string]string)
//val := reflect.ValueOf(t).Elem()
//for i := 0; i < val.NumField(); i++ {
//field := val.Type().Field(i)
//nameTag := field.Tag.Get("name")
//descriptionTag := field.Tag.Get("description")
//if nameTag != "" {
//fieldMeta[field.Name] = map[string]string{
//"name":        nameTag,
//"description": descriptionTag,
//}
//}
//}
//return fieldMeta
//}

//func GetMetadataStrategyWavesSettingsMetadata() string {
//meta := getFieldMetadata(&WawesSettings{})
//metaJson, err := json.Marshal(meta)
//if err != nil {
//return fmt.Errorf("failed to marshal metadata: %w", err).Error()
//}
//return string(metaJson)
//}

// // UnmarshalConfig распаковывает конфигурацию стратегии
// func (s *Strategy) UnmarshalConfig() (map[string]interface{}, error) {
// 	var config map[string]interface{}
// 	if err := json.Unmarshal(s.Config, &config); err != nil {
// 		return nil, err
// 	}
// 	return config, nil
// }

// // MarshalConfig упаковывает конфигурацию стратегии
// func (s *Strategy) MarshalConfig(config map[string]interface{}) error {
// 	configJSON, err := json.Marshal(config)
// 	if err != nil {
// 		return err
// 	}
// 	s.Config = configJSON
// 	return nil
// }
