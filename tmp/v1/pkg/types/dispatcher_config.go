package types

type SignalType string

type DispatcherConfig struct {
	Rules    []RuleConfig                   `json:"rules"`
	Handlers map[SignalType][]HandlerConfig `json:"handlers"`
}

type RuleConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type HandlerConfig struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

func (c DispatcherConfig) GetComponentType() string {
	return "dispatcher"
}

func (c RuleConfig) GetComponentType() string {
	return "rule"
}

func (c HandlerConfig) GetComponentType() string {
	return "handler"
}
