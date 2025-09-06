package types

type ComponentRepository interface {
	GetByID(id string) (*Component, error)
	GetByIDs(ids []string) ([]Component, error)
}
