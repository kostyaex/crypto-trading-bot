// Обновление текущего времени
package systems

import (
	"crypto-trading-bot/internal/engine/resources"
	"fmt"
	"time"

	"github.com/andygeiss/ecs"
)

type timeUpdateSystem struct {
	CurrentTime *resources.CurrentTime
}

func (s *timeUpdateSystem) Process(em ecs.EntityManager) (state int) {

	s.CurrentTime.Timestamp = time.Now().Unix()
	fmt.Printf("SET CURRENT TIME: %s\n", time.Unix(s.CurrentTime.Timestamp, 0).Format("2006-01-02 15:04:05"))

	return ecs.StateEngineContinue
}

func (s *timeUpdateSystem) Setup() {}

func (s *timeUpdateSystem) Teardown() {}

func NewTimeUpdateSystem(CurrentTime *resources.CurrentTime) ecs.System {
	return &timeUpdateSystem{
		CurrentTime: CurrentTime,
	}
}
