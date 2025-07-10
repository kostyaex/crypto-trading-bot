package series

import (
	"crypto-trading-bot/internal/models"
	"fmt"
	"time"
)

type Point struct {
	Value      float64
	Weight     float64
	Time       time.Time
	MarketData *models.MarketData // основание по которому формировалась точка
}

type Series struct {
	Points []Point
}

// Возвращает первую точку
func (s *Series) First() *Point {
	if len(s.Points) == 0 {
		return nil
	}
	return &s.Points[0]
}

// Возвращает последнюю точку
func (s *Series) Last() *Point {
	if len(s.Points) == 0 {
		return nil
	}
	return &s.Points[len(s.Points)-1]
}

// Формирует строку для вывода в лог
func (s *Series) String() string {
	if len(s.Points) == 0 {
		return "[]"
	}

	if len(s.Points) == 1 {
		first := s.Points[0].MarketData
		return fmt.Sprintf("[%s $%.2f]", first.Timestamp.Format(time.DateTime), first.ClusterPrice)
	}

	first := s.Points[0].MarketData
	last := s.Points[len(s.Points)-1].MarketData

	return fmt.Sprintf("(%d)[%s $%.2f - %s $%.2f]",
		len(s.Points),
		first.Timestamp.Format(time.DateTime), first.ClusterPrice,
		last.Timestamp.Format(time.DateTime), last.ClusterPrice)
}

type SeriesBuilder interface {
	AddPoints(activeSeries []Series, newPoints []Point) []Series
}

type BuilderType string

const (
	SimpleBuilderType   BuilderType = "simple"
	WindowedBuilderType BuilderType = "windowed"
)

func NewSeriesBuilder(config map[string]interface{}) (SeriesBuilder, error) {
	builderType, ok := config["type"].(string)
	if !ok {
		return nil, fmt.Errorf("не указан тип стратегии")
	}

	switch BuilderType(builderType) {
	case SimpleBuilderType:
		return &SimpleSeriesBuilder{
			valueFactor: config["value_factor"].(float64),
			timeFactor:  config["time_factor"].(float64),
		}, nil
	case WindowedBuilderType:
		return &WindowedSeriesBuilder{
			maxTimeGap:   time.Duration(config["max_time_gap"].(int)) * time.Second,
			maxValueDiff: config["max_value_diff"].(float64),
		}, nil
	default:
		return nil, fmt.Errorf("неизвестный тип стратегии: %s", builderType)
	}
}
