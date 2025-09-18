package ecsx

import "github.com/andygeiss/ecs"

type customEngine struct {
	ecs.Engine
	state         int
	entityManager EntityManager
	systemManager ecs.SystemManager
}

func NewCustomEngine(entityManager EntityManager, systemManager ecs.SystemManager) ecs.Engine {
	return &customEngine{
		Engine:        ecs.NewDefaultEngine(entityManager, systemManager),
		entityManager: entityManager,
		systemManager: systemManager,
	}
}

// Run calls the Process() method for each System
// until ShouldEngineStop is set to true.
func (e *customEngine) Run() {
	shouldStop := false
	for !shouldStop {
		for _, system := range e.systemManager.Systems() {
			e.state = system.Process(e.entityManager)
			if e.state == ecs.StateEngineStop {
				shouldStop = true
				break
			}
		}
		e.entityManager.ProcessEvents()
	}
}
