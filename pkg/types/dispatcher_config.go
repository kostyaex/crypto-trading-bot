package types

type SignalType string

type DispatcherSettings struct {
	Rules    []RuleSettings                   `json:"rules"`
	Handlers map[SignalType][]HandlerSettings `json:"handlers"`
}

type RuleSettings struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

type HandlerSettings struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config"`
}

func (c DispatcherSettings) GetComponentType() string {
	return "dispatcher"
}
