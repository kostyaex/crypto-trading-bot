package series

import (
	"crypto-trading-bot/pkg/types"
	"fmt"
	"sync"
	"time"
)

type SeriesBuilder struct {
	algorithm Algorithm
	series    []*Series
	mutex     sync.RWMutex
	timestamp time.Time // время обновления последних данных
}

func NewSeriesBuilder(config map[string]interface{}) (*SeriesBuilder, error) {
	builderType, ok := config["type"].(string)
	if !ok {
		return nil, fmt.Errorf("не указан тип стратегии")
	}

	builder := &SeriesBuilder{
		series: make([]*Series, 0),
	}

	switch AlgorithmType(builderType) {
	case SimpleAlgorithmType:
		builder.algorithm = &SimpleAlgorithm{
			valueFactor: config["value_factor"].(float64),
			timeFactor:  config["time_factor"].(float64),
		}
		return builder, nil
	case WindowedAlgorithmType:
		builder.algorithm = &WindowedAlgorithm{
			maxTimeGap:   time.Duration(config["max_time_gap"].(int)) * time.Second,
			maxValueDiff: config["max_value_diff"].(float64),
		}
		return builder, nil
	default:
		return nil, fmt.Errorf("неизвестный тип стратегии: %s", builderType)
	}
}

func (builder *SeriesBuilder) AddClusteredData(newData []*types.MarketData) {
	builder.mutex.Lock()
	defer builder.mutex.Unlock()

	var points []Point
	for _, md := range newData {
		point := Point{
			Value:      md.ClusterPrice,
			Weight:     md.Volume,
			Time:       md.Timestamp,
			MarketData: md,
		}
		points = append(points, point)
	}

	builder.series = builder.algorithm.AddPoints(builder.series, points)
	builder.timestamp = newData[len(newData)-1].Timestamp
}

// возвращает все серии
func (builder *SeriesBuilder) GetSeries() []*Series {
	return builder.series
}

// Возвращает обновлённые серии
func (builder *SeriesBuilder) GetActiveSeries() []*Series {
	builder.mutex.RLock()
	defer builder.mutex.RUnlock()

	activeSeries := make([]*Series, 0)

	for _, sr := range builder.series {
		// проверяем последнюю точку серии
		last := sr.Last()
		if last == nil || last.Time.Before(builder.timestamp) {
			continue
		}

		activeSeries = append(activeSeries, sr)
	}

	return activeSeries
}
