package stages

import (
	"crypto-trading-bot/internal/components"
	"crypto-trading-bot/internal/types"
)

type HistoricalSource struct {
	settings components.HistoricalSourceSettings
}

func (s *HistoricalSource) UpdateConfig(comps ...types.Component) {
	for _, c := range comps {
		if val, ok := c.(components.HistoricalSourceSettings); ok {
			s.settings = val
		}
	}
}
