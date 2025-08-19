package series

import (
	"math"
	"time"
)

type WindowedAlgorithm struct {
	maxTimeGap   time.Duration
	maxValueDiff float64
}

func (a *WindowedAlgorithm) Name() string {
	return "WindowedAlgorithm"
}

// на основе временного окна и кластеризации значений
func (b *WindowedAlgorithm) AddPoints(activeSeries []*Series, newPoints []Point) []*Series {
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
			activeSeries = append(activeSeries, &Series{Points: []Point{pt}})
		}
	}
	return activeSeries
}
