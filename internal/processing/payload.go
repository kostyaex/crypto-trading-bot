package processing

import (
	"crypto-trading-bot/pkg/pipeline"
	"sync"
	"time"
)

var (
	// проверка соответствия интерфейсу;
	_ pipeline.Payload = (*TradingPayload)(nil)

	PayloadPool = sync.Pool{
		New: func() interface{} { return new(TradingPayload) },
	}
)

type TradingPayload struct {
	// Этап источника получения торговых данных
	Symbol    string
	Interval  string
	StartTime time.Time
	EndTime   time.Time
	// Этап выборки данных
	CurrentPrice float64
}

// Реализация интерфейса pipeline.Payload

// Clone implements pipeline.Payload.
func (p *TradingPayload) Clone() pipeline.Payload {
	newP := PayloadPool.Get().(*TradingPayload)
	newP.CurrentPrice = p.CurrentPrice

	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *TradingPayload) MarkAsProcessed() {
	// Очистка
	p.Symbol = ""
	p.Interval = ""
	p.StartTime = time.Time{}
	p.EndTime = time.Time{}
	p.CurrentPrice = 0
	PayloadPool.Put(p)
}
