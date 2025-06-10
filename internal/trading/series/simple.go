package series

import "math"

type SimpleSeriesBuilder struct {
	valueFactor float64
	timeFactor  float64
}

// Формирование серий по сопоставлению последней точки серии с добавляемыми
func (b *SimpleSeriesBuilder) AddPoints(activeSeries []Series, newPoints []Point) []Series {
	for _, pt := range newPoints {
		bestMatch := -1
		minScore := math.MaxFloat64

		for i := range activeSeries {
			lastPt := activeSeries[i].Points[len(activeSeries[i].Points)-1]
			dt := pt.Time.Sub(lastPt.Time).Seconds()
			dv := math.Abs(pt.Value - lastPt.Value)
			combinedWeight := pt.Weight + lastPt.Weight + 1e-9

			score := dv/b.valueFactor/combinedWeight + dt/b.timeFactor/combinedWeight

			if dt >= 0 && score < minScore {
				minScore = score
				bestMatch = i
			}
		}

		if bestMatch != -1 && minScore < 1.0 {
			activeSeries[bestMatch].Points = append(activeSeries[bestMatch].Points, pt)
		} else {
			activeSeries = append(activeSeries, Series{Points: []Point{pt}})
		}
	}
	return activeSeries
}
