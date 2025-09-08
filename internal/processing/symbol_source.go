package processing

import (
	"context"
	"crypto-trading-bot/pkg/pipeline"
)

type SymbolSource struct {
	items []SymbolItem
	index int
}

type SymbolItem struct {
	Symbol   string
	Interval string
	//Exchange string
}

func NewSymbolSource(items []SymbolItem) *SymbolSource {
	return &SymbolSource{
		items: items,
		index: -1,
	}
}

func (s *SymbolSource) Error() error { return nil }

func (s *SymbolSource) Next(context.Context) bool {
	if s.index == len(s.items)-1 {
		return false
	}

	//fmt.Printf("data: %v\n", s.data[s.index])

	s.index++
	return true
}

func (s *SymbolSource) Payload() pipeline.Payload {
	p := PayloadPool.Get().(*TradingPayload)

	p.Symbol = s.items[s.index].Symbol
	p.Interval = s.items[s.index].Interval

	return p
}
