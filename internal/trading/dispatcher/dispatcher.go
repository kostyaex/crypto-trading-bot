package dispatcher

import (
	"crypto-trading-bot/internal/service/series"
	"time"
)

type TradeSignal struct {
	Timestamp   time.Time
	Symbol      string
	Price       float64
	Volume      float64
	BuyVolume   float64
	SellVolume  float64
	Type        SignalType
	Series      *series.Series
	Description string
}

type SignalType string

const (
	SignalBuy  SignalType = "buy"
	SignalSell SignalType = "sell"
	SignalHold SignalType = "hold"
)

type SignalRule interface {
	Evaluate(series *series.Series) (TradeSignal, bool)
}

type ActionHandler interface {
	Handle(signal TradeSignal)
}

type ActionHandlerFunc func(signal TradeSignal)

type Dispatcher struct {
	rules    []SignalRule
	handlers map[SignalType][]ActionHandler
}

func NewDispatcher(rules ...SignalRule) *Dispatcher {
	return &Dispatcher{
		rules:    rules,
		handlers: make(map[SignalType][]ActionHandler),
	}
}

// Регистрация обработчика
func (d *Dispatcher) Register(signalType SignalType, handler ActionHandler) {
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
