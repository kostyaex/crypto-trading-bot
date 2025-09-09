package settings

type SettingsRepository interface {
	GetByID(id string) (*Settings, error)
	GetByIDs(ids []string) ([]Settings, error)
}
