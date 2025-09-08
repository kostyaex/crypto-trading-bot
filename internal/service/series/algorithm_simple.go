package series

import (
	"crypto-trading-bot/internal/types"
	"math"
)

type SimpleAlgorithm struct {
	valueFactor float64
	timeFactor  float64
}

func (a *SimpleAlgorithm) Name() string {
	return "SimpleAlgorithm"
}

// Формирование серий по сопоставлению последней точки серии с добавляемыми
func (b *SimpleAlgorithm) AddPoints(activeSeries []*types.Series, newPoints []types.Point) []*types.Series {

	// Точки с одинаковым временем, не должны попадать в одну серию

	for _, pt := range newPoints {
		bestMatch := -1
		minScore := math.MaxFloat64

		for i := range activeSeries {

			// Получение последней точки серии
			lastPt := activeSeries[i].Points[len(activeSeries[i].Points)-1]

			if pt.Time == lastPt.Time {
				continue
			}

			dt := pt.Time.Sub(lastPt.Time).Seconds()
			dv := math.Abs(pt.Value - lastPt.Value)
			combinedWeight := pt.Weight + lastPt.Weight + 1e-9

			score := dv*b.valueFactor/combinedWeight + dt*b.timeFactor/combinedWeight

			if dt >= 0 && score < minScore {
				minScore = score
				bestMatch = i
			}
		}

		if bestMatch != -1 && minScore < 1.0 {
			activeSeries[bestMatch].Points = append(activeSeries[bestMatch].Points, pt)
		} else {
			activeSeries = append(activeSeries, &types.Series{Points: []types.Point{pt}})
		}
	}
	return activeSeries
}
