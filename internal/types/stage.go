package types

type Stage struct {
	Name       string   `json:"name"`
	Components []string `json:"components"` // список ID или имён компонентов, которые нужны на этой стадии
}

type StageConfig struct {
	StageName  string
	Components map[string]Component // мапа: component_id -> config
}
