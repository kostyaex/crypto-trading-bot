package types

import "context"

type Runnable interface {
	Run(ctx context.Context) error
}

type Processor interface {
	Process(ctx context.Context, data interface{}) (interface{}, error)
}
