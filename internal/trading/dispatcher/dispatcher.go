package dispatcher

import (
	"crypto-trading-bot/internal/service/series"
	"crypto-trading-bot/pkg/types"
	"time"
)

type TradeSignal struct {
	Timestamp   time.Time
	Symbol      string
	Price       float64
	Volume      float64
	BuyVolume   float64
	SellVolume  float64
	Type        types.SignalType
	Series      *series.Series
	Description string
}

type SignalRule interface {
	Evaluate(series *series.Series) (TradeSignal, bool)
}

type ActionHandler interface {
	Handle(signal TradeSignal)
}

type ActionHandlerFunc func(signal TradeSignal)

type Dispatcher struct {
	rules    []SignalRule
	handlers map[types.SignalType][]ActionHandler
}

func NewDispatcher(rules ...SignalRule) *Dispatcher {
	return &Dispatcher{
		rules:    rules,
		handlers: make(map[types.SignalType][]ActionHandler),
	}
}

// Регистрация обработчика
func (d *Dispatcher) Register(signalType types.SignalType, handler ActionHandler) {
	d.handlers[signalType] = append(d.handlers[signalType], handler)
}

func (d *Dispatcher) Dispatch(series *series.Series) {

	resultSignal := TradeSignal{Type: SignalHold}

	for _, rule := range d.rules {
		if signal, ok := rule.Evaluate(series); ok {
			resultSignal = signal
			break
		}
	}

	for _, handler := range d.handlers[resultSignal.Type] {
		handler.Handle(resultSignal)
	}

}
