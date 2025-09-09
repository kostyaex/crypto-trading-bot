package ecsx

import "github.com/andygeiss/ecs"

type customEngine struct {
	entityManager ecs.EntityManager
	systemManager ecs.SystemManager
	state         int
}

//kostyaex: добавил state и метод IsDone

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
	}
}

// Tick calls the Process() method for each System exactly once
func (e *customEngine) Tick() {
	for _, system := range e.systemManager.Systems() {
		e.state = system.Process(e.entityManager)
		if e.state == ecs.StateEngineStop {
			break
		}
	}
}

// Setup calls the Setup() method for each System
// and initializes ShouldEngineStop and ShouldEnginePause with false.
func (e *customEngine) Setup() {
	for _, sys := range e.systemManager.Systems() {
		sys.Setup()
	}
}

// Teardown calls the Teardown() method for each System.
func (e *customEngine) Teardown() {
	for _, sys := range e.systemManager.Systems() {
		sys.Teardown()
	}
}

func (e *customEngine) IsDone() bool {
	if e.state == ecs.StateEngineStop {
		return true
	}
	return false
}

func NewCustomEngine(entityManager ecs.EntityManager, systemManager ecs.SystemManager) Engine {
	return &customEngine{
		entityManager: entityManager,
		systemManager: systemManager,
	}
}
