package settings

type Settings interface {
	// ComponentType возвращает уникальный строковый идентификатор типа компонента.
	// Должен совпадать с ID, по которому компонент регистрируется в реестре.
	SettingsType() string
}

// type ComponentConfig struct {
// 	ID       string          `db:"id" json:"id"`
// 	Name     string          `db:"name" json:"name"`
// 	Settings json.RawMessage `db:"settings" json:"settings"`
// }
