package pipelinemanager

import (
	"crypto-trading-bot/pkg/pipeline"
	"sync"
)

var (
	// проверка соответствия интерфейсу;
	_ pipeline.Payload = (*tradingPayload)(nil)

	payloadPool = sync.Pool{
		New: func() interface{} { return new(tradingPayload) },
	}
)

type tradingPayload struct {
	CurrentPrice int64
}

// Clone implements pipeline.Payload.
func (p *tradingPayload) Clone() pipeline.Payload {
	newP := payloadPool.Get().(*tradingPayload)
	newP.CurrentPrice = p.CurrentPrice

	return newP
}

// MarkAsProcessed implements pipeline.Payload
func (p *tradingPayload) MarkAsProcessed() {
	p.CurrentPrice = 0
	payloadPool.Put(p)
}
