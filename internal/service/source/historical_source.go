package source

import (
	"crypto-trading-bot/internal/components"
	"crypto-trading-bot/internal/service/marketdata"
	"crypto-trading-bot/internal/types"
	"fmt"
)

type HistoricalSource struct {
	settings          components.HistoricalSourceSettings
	marketDataService *marketdata.MarketDataService
	data              []*types.MarketData
	index             int
}

func NewHistoricalSource(marketDataService marketdata.MarketDataService, comps ...types.Component) (*HistoricalSource, error) {

	var err error

	s := &HistoricalSource{
		marketDataService: &marketDataService,
	}

	s.UpdateConfig(comps...)

	s.data, err = marketDataService.GetMarketDataPeriod(
		s.settings.Symbol,
		s.settings.Interval,
		s.settings.StartTime,
		s.settings.EndTime,
	)

	return s, err
}

func (s *HistoricalSource) UpdateConfig(comps ...types.Component) {
	for _, c := range comps {
		if val, ok := c.(*components.HistoricalSourceSettings); ok {
			s.settings = *val
		}
	}
}

func (s *HistoricalSource) Next() bool {
	//fmt.Printf("data: %d\n", len(s.data))

	if s.index == len(s.data) {
		return false
	}

	fmt.Printf("data: %v\n", s.data[s.index])

	s.index++
	return true
}
