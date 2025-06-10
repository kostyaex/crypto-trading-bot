package dispatcher

import (
	"crypto-trading-bot/internal/trading/series"
	"time"
)

type TradeSignal struct {
	Timestamp  time.Time
	Symbol     string
	Price      float64
	Volume     float64
	BuyVolume  float64
	SellVolume float64
	Type       SignalType
	Series     *series.Series
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

type SignalDispatcher struct {
	rules []SignalRule
}

func NewSignalDispatcher(rules ...SignalRule) *SignalDispatcher {
	return &SignalDispatcher{rules: rules}
}

func (d *SignalDispatcher) Dispatch(series *series.Series) (TradeSignal, bool) {
	for _, rule := range d.rules {
		if signal, ok := rule.Evaluate(series); ok {
			return signal, true
		}
	}
	return TradeSignal{Type: SignalHold}, false
}
