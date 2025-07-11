package series

import (
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

	var series []Series

	now := time.Now()

	// Итерация 1
	points1 := []Point{
		{Value: 100.0, Weight: 2.5, Time: now},
		{Value: 101.0, Weight: 3.0, Time: now.Add(1 * time.Second)},
	}
	series = builder.AddPoints(series, points1)

	// Итерация 2
	points2 := []Point{
		{Value: 102.0, Weight: 4.0, Time: now.Add(2 * time.Second)},
		{Value: 200.0, Weight: 1.0, Time: now.Add(2 * time.Second)},
	}
	series = builder.AddPoints(series, points2)

	// Вывод
	for i, s := range series {
		fmt.Printf("Серия %d:\n", i+1)
		for _, p := range s.Points {
			fmt.Printf("  Значение: %.2f | Вес: %.2f | Время: %s\n", p.Value, p.Weight, p.Time.Format("15:04:05"))
		}
	}

	// Сбор метрик
	metrics := CollectMetrics(series)
	metrics.Print()

}
