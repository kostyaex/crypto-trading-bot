package series

import (
	"math"
	"time"
)

type WindowedSeriesBuilder struct {
	maxTimeGap   time.Duration
	maxValueDiff float64
}

// на основе временного окна и кластеризации значений
func (b *WindowedSeriesBuilder) AddPoints(activeSeries []Series, newPoints []Point) []Series {
	for _, pt := range newPoints {
		matched := false
		for i := range activeSeries {
			lastPt := activeSeries[i].Points[len(activeSeries[i].Points)-1]
			if pt.Time.Sub(lastPt.Time) <= b.maxTimeGap && math.Abs(pt.Value-lastPt.Value) < b.maxValueDiff {
				activeSeries[i].Points = append(activeSeries[i].Points, pt)
				matched = true
				break
			}
		}
		if !matched {
			activeSeries = append(activeSeries, Series{Points: []Point{pt}})
		}
	}
	return activeSeries
}
