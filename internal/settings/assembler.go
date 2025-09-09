package settings

// type ConfigAssembler struct {
// 	repo     types.ComponentRepository
// 	registry *ComponentRegistry
// }

// func NewConfigAssembler(repo types.ComponentRepository, registry *ComponentRegistry) *ConfigAssembler {
// 	return &ConfigAssembler{repo: repo, registry: registry}
// }

//func (a *ConfigAssembler) BuildStageConfig(stage types.Stage) (*types.StageConfig, error) {
//components, err := a.repo.GetByIDs(stage.Components)
//if err != nil {
//return nil, fmt.Errorf("failed to fetch components for stage %s: %w", stage.Name, err)
//}

//componentMap := make(map[string]types.Component)

//for _, comp := range components {
//typedComponent, err := a.registry.Build(comp.ID, comp.Settings)
//if err != nil {
//return nil, fmt.Errorf("failed to build component %s: %w", comp.ID, err)
//}
//componentMap[comp.ID] = typedComponent
//}

//// Проверка: все ли компоненты найдены и собраны?
//for _, id := range stage.Components {
//if _, exists := componentMap[id]; !exists {
//return nil, fmt.Errorf("component %s not built (should not happen)", id)
//}
//}

//return &StageConfig{
//StageName:  stage.Name,
//Components: componentMap,
//}, nil
//}
