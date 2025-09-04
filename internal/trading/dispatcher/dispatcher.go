package dispatcher

import (
	"context"
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
	Series      *types.Series
	Description string
}

type SignalRule interface {
	Evaluate(series *types.Series) (TradeSignal, bool)
}

type ActionHandler interface {
	Handle(signal TradeSignal)
}

type ActionHandlerFunc func(signal TradeSignal)

type Dispatcher struct {
	rules    []SignalRule
	handlers map[types.SignalType][]ActionHandler
}

// Регистрация обработчика
func (d *Dispatcher) Register(signalType types.SignalType, handler ActionHandler) {
	d.handlers[signalType] = append(d.handlers[signalType], handler)
}

func (d *Dispatcher) Dispatch(ctx context.Context, series *types.Series) {

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
