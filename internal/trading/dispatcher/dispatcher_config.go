package dispatcher

import (
	"crypto-trading-bot/internal/engine"
	"crypto-trading-bot/pkg/types"
	"fmt"
	"log"
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

func NewDispatcher(config *types.DispatcherConfig) (*Dispatcher, error) {

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
	dispatcher := &Dispatcher{
		rules:    rules,
		handlers: make(map[types.SignalType][]ActionHandler),
	}

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

// перенести в файл factory.go
func NewDispatcherFactory() engine.ComponentFactory {
	return func(config types.ComponentConfig, logger *log.Logger) (*engine.Component, error) {
		typedConfig, ok := config.(*types.DispatcherConfig)
		if !ok {
			return nil, fmt.Errorf("invalid config type")
		}

		var comp engine.Component

		dispatcher, err := NewDispatcher(typedConfig)

		if err != nil {
			return nil, fmt.Errorf("invalid config type")
		}

		comp.Dispatcher = dispatcher

		return &comp, nil
	}
}

// как выглядил пример в решении ИИ:
// func NewTradingFactory(tradeStore storage.TradeStore) engine.ComponentFactory {
// 	return func(config types.ComponentConfig, logger *log.Logger) (*engine.Component, error) {
// 		typedConfig, ok := config.(*types.TradingConfig)
// 		if !ok {
// 			return nil, fmt.Errorf("invalid config type")
// 		}

// 		var comp engine.Component

// 		switch typedConfig.Mode {
// 		case "live":
// 			executor := NewLiveExecutor(typedConfig, logger, tradeStore)
// 			comp.Processor = executor
// 		case "paper":
// 			executor := NewPaperExecutor(typedConfig, logger, tradeStore)
// 			comp.Processor = executor
// 		case "backtest":
// 			executor := NewBacktestExecutor(typedConfig, logger, tradeStore)
// 			comp.Runnable = executor // может быть Runnable, если генерирует сигналы
// 		default:
// 			return nil, fmt.Errorf("unsupported mode: %s", typedConfig.Mode)
// 		}

// 		return &comp, nil
// 	}
// }
