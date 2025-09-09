package settings

type ConfigUpdate interface {
	UpdateConfig(...Settings)
}
