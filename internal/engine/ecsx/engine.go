package ecsx

import "github.com/andygeiss/ecs"

type Engine interface {
	ecs.Engine
	IsDone() bool
}
