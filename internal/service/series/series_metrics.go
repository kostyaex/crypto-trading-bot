package series

import (
	"crypto-trading-bot/internal/types"
	"fmt"
	"math"
	"time"
)

type SeriesMetrics struct {
	TotalSeries        int       // Общее число серий
	TotalPoints        int       // Общее число точек во всех сериях
	AvgLength          float64   // Средняя длина серии
	MaxLength          int       // Максимальная длина серии
	MinLength          int       // Минимальная длина серии
	AvgValue           float64   // Среднее значение по всем точкам
	MinValue, MaxValue float64   // Мин/макс значения
	StartTime, EndTime time.Time // Самая ранняя и поздняя точки
}

func CollectMetrics(series []*types.Series) SeriesMetrics {
	var m SeriesMetrics
	m.TotalSeries = len(series)

	if m.TotalSeries == 0 {
		return m
	}

	var totalValuesSum float64
	m.MinValue = math.MaxFloat64
	m.MaxValue = -math.MaxFloat64
	m.StartTime = time.Now().Add(24 * time.Hour)           // будем искать минимум
	m.EndTime, _ = time.Parse(time.DateOnly, "2000-01-01") // будем искать максимум

	for _, s := range series {
		length := len(s.Points)
		m.TotalPoints += length
		if length > m.MaxLength {
			m.MaxLength = length
		}
		if length < m.MinLength || m.MinLength == 0 {
			m.MinLength = length
		}

		for _, p := range s.Points {
			totalValuesSum += p.Value
			if p.Value < m.MinValue {
				m.MinValue = p.Value
			}
			if p.Value > m.MaxValue {
				m.MaxValue = p.Value
			}
			if p.Time.Before(m.StartTime) {
				m.StartTime = p.Time
			}
			if p.Time.After(m.EndTime) {
				m.EndTime = p.Time
			}
		}
	}

	if m.TotalPoints > 0 {
		m.AvgValue = totalValuesSum / float64(m.TotalPoints)
	}
	m.AvgLength = float64(m.TotalPoints) / float64(m.TotalSeries)

	return m
}

func (metrics SeriesMetrics) Print() {
	fmt.Printf("Число серий: %d\n", metrics.TotalSeries)
	fmt.Printf("Общее число точек: %d\n", metrics.TotalPoints)
	fmt.Printf("Средняя длина серии: %.2f\n", metrics.AvgLength)
	fmt.Printf("Мин/макс длина серии: %d / %d\n", metrics.MinLength, metrics.MaxLength)
	fmt.Printf("Среднее значение: %.2f\n", metrics.AvgValue)
	fmt.Printf("Мин/макс значение: %.2f / %.2f\n", metrics.MinValue, metrics.MaxValue)
	fmt.Printf("Период данных: с %s до %s\n", metrics.StartTime.Format(time.RFC3339), metrics.EndTime.Format(time.RFC3339))
}
