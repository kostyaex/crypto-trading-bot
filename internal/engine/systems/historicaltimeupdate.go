// Обновление текущего времени
package systems

import (
	"crypto-trading-bot/internal/engine/components"
	"crypto-trading-bot/internal/engine/resources"
	"fmt"
	"time"

	"github.com/andygeiss/ecs"
)

type HistoricalTimeUpdateSystem struct {
	CurrentTime *resources.CurrentTime
	Candles     []components.Candle
	Index       int
}

func (s *HistoricalTimeUpdateSystem) Process(em ecs.EntityManager) (state int) {

	if s.Index >= len(s.Candles) {
		return ecs.StateEngineStop
	}

	// e := em.Get("currenttime")
	// currentTime := e.Get(components.MaskCurrentTime).(*components.CurrentTime)
	// currentTime.Timestamp = s.Candles[s.Index].Timestamp

	s.CurrentTime.Timestamp = s.Candles[s.Index].Timestamp

	fmt.Printf("SET CURRENT TIME: %s\n", time.Unix(s.CurrentTime.Timestamp, 0).Format("2006-01-02 15:04:05"))

	s.Index++

	return ecs.StateEngineContinue
}

func (s *HistoricalTimeUpdateSystem) Setup() {}

func (s *HistoricalTimeUpdateSystem) Teardown() {}

func NewHistoricalTimeUpdateSystem(CurrentTime *resources.CurrentTime, Candles []components.Candle) ecs.System {
	return &HistoricalTimeUpdateSystem{
		CurrentTime: CurrentTime,
		Candles:     Candles,
	}
}
