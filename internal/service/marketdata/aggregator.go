package marketdata

import (
	"crypto-trading-bot/internal/service/clusters"
	"crypto-trading-bot/internal/types"
)

// Сгруппировать торговые данные по количеству в блоке
type Aggregator struct {
	BlockSize int
	block     []*types.MarketData
}

func NewAggregator(blockSize int) *Aggregator {
	return &Aggregator{
		BlockSize: blockSize,
		block:     make([]*types.MarketData, 0),
	}
}

func (a *Aggregator) Add(data *types.MarketData) {
	if len(a.block) == a.BlockSize {
		a.block = a.block[:0] // Очищаем блок
	}
	a.block = append(a.block, data)
}

func (a *Aggregator) IsReady() bool {
	if len(a.block) == a.BlockSize {
		return true
	} else {
		return false
	}
}

func (a *Aggregator) GetAggregatedData() []*types.MarketData {
	if len(a.block) == a.BlockSize {
		return a.block
	} else {
		return nil
	}
}

func (a *Aggregator) GetClusteredData(NumClusters int, setTimeframe string) []*types.MarketData {
	if len(a.block) != a.BlockSize {
		return nil
	}

	// здесь сворачиваем данные в кластеры. Т.е. к примеру данные за секундный интервал в 5 минутный, получим столько значений, сколько указано количество кластеров.
	clusteredMd := clusters.ClusterMarketData(a.block, setTimeframe, NumClusters)

	return clusteredMd

}
