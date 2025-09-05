package types

import (
	"fmt"
	"time"
)

type Point struct {
	Value      float64     `json:"value"`
	Weight     float64     `json:"weight"`
	Time       time.Time   `json:"time"`
	MarketData *MarketData `json:"market_data"` // основание по которому формировалась точка
}

type Series struct {
	Points []Point `json:"points"`
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
