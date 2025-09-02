package dispatcher

import (
	"crypto-trading-bot/pkg/types"
	"encoding/json"
	"fmt"
)

const (
	SignalBuy  types.SignalType = "buy"
	SignalSell types.SignalType = "sell"
	SignalHold types.SignalType = "hold"
)

var (
	ruleRegistry    = make(map[string]func(map[string]interface{}) (SignalRule, error))
	handlerRegistry = make(map[string]func(map[string]interface{}) (ActionHandler, error))
)

func init() {
	ruleRegistry["volume_trand"] = NewVolumeTrendRule
	handlerRegistry["logger"] = NewLoggerHandler
	handlerRegistry["file"] = NewFileLoggerHandler
}

func NewDispatcherFromJSON(jsonData []byte) (*Dispatcher, error) {
	var config types.DispatcherSettings

	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	// Создаём правила
	var rules []SignalRule
	for _, r := range config.Rules {
		factory, ok := ruleRegistry[r.Type]
		if !ok {
			return nil, fmt.Errorf("unknown rule type: %s", r.Type)
		}
		rule, err := factory(r.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule %s: %v", r.Type, err)
		}
		rules = append(rules, rule)
	}

	// Создаём диспетчер
	dispatcher := NewDispatcher(rules...)

	// Регистрируем обработчики
	for signalType, handlers := range config.Handlers {
		for _, h := range handlers {
			factory, ok := handlerRegistry[h.Type]
			if !ok {
				return nil, fmt.Errorf("unknown handler type: %s", h.Type)
			}
			handler, err := factory(h.Config)
			if err != nil {
				return nil, fmt.Errorf("failed to create handler %s: %v", h.Type, err)
			}
			dispatcher.Register(signalType, handler)
		}
	}

	return dispatcher, nil
}
