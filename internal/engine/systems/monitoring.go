// Отслеживание состояния бота
package systems

import (
	"crypto-trading-bot/internal/engine/resources"
	"fmt"
	"time"

	"github.com/andygeiss/ecs"
)

type monitoringSystem struct {
	CurrentTime *resources.CurrentTime
}

func (s *monitoringSystem) Process(em ecs.EntityManager) (state int) {
	fmt.Printf("MON CURRENT TIME: %s\n", time.Unix(s.CurrentTime.Timestamp, 0).Format("2006-01-02 15:04:05"))
	return ecs.StateEngineContinue
}

func (s *monitoringSystem) Setup() {}

func (s *monitoringSystem) Teardown() {}

func NewMonitoringSystem(currentTime *resources.CurrentTime) ecs.System {
	return &monitoringSystem{
		CurrentTime: currentTime,
	}
}
