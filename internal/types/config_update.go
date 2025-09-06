package types

type ConfigUpdate interface {
	UpdateConfig(...Component)
}
