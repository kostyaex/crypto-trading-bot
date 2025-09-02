package types

import (
	"context"
)

type Runnable interface {
	Run(ctx context.Context) error
}

type Processor interface {
	Process(ctx context.Context, data interface{}) (interface{}, error)
}

type MarketDataSource interface {
	Run(ctx context.Context, output chan<- *MarketData) error
}

// type SignalGenerator interface {
// 	Run(ctx context.Context, input <-chan *Candle, output chan<- *Signal) error
// }

// type RiskProcessor interface {
// 	Process(ctx context.Context, signal *Signal) (*Signal, error)
// }

// type OrderExecutor interface {
// 	Process(ctx context.Context, signal *Signal) error
// }
