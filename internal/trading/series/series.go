package series

import (
	"fmt"
	"time"
)

type Point struct {
	Value  float64
	Weight float64
	Time   time.Time
}

type Series struct {
	Points []Point
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
