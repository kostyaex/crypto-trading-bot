package series

import (
	"crypto-trading-bot/pkg/types"
	"fmt"
	"testing"
	"time"
)

func TestNewSeriesBuilder(t *testing.T) {

	// Пример конфигурации
	config := map[string]interface{}{
		"type":         "simple",
		"value_factor": 1.0,
		"time_factor":  100.0,
	}

	// Создаем стратегию
	builder, err := NewSeriesBuilder(config)
	if err != nil {
		panic(err)
	}

	var series []*types.Series

	now := time.Now()

	// Итерация 1
	// Два значения с одинаковым временем, не должны попасть в одну серию
	points1 := []types.Point{
		{Value: 100.0, Weight: 2.5, Time: now},
		{Value: 101.0, Weight: 3.0, Time: now},
		{Value: 110.0, Weight: 2.0, Time: now},
	}
	series = builder.algorithm.AddPoints(series, points1)

	// Итерация 2
	points2 := []types.Point{
		{Value: 102.0, Weight: 4.0, Time: now.Add(2 * time.Second)},
		{Value: 200.0, Weight: 1.0, Time: now.Add(2 * time.Second)},
	}
	series = builder.algorithm.AddPoints(series, points2)

	// Вывод
	for i, s := range series {
		fmt.Printf("Серия %d:\n", i+1)

		// Проверяем, чтобы в серии не было точек с одинаковым временем
		_time := time.Time{}
		for _, p := range s.Points {
			if p.Time.Equal(_time) {
				t.Errorf("В серии %d есть точки с одинаковым временем: %s", i+1, p.Time.Format("15:04:05"))
			} else {
				_time = p.Time
			}
			fmt.Printf("  Значение: %.2f | Вес: %.2f | Время: %s\n", p.Value, p.Weight, p.Time.Format("15:04:05"))
		}
	}

	// Сбор метрик
	metrics := CollectMetrics(series)
	metrics.Print()

}
